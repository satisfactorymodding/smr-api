package validation

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/Vilsol/slox"
	"github.com/VirusTotal/vt-go"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var client *vt.Client

func InitializeVirusTotal() {
	client = vt.NewClient(viper.GetString("virustotal.key"))

	if client == nil {
		panic("failed to initialize virustotal client")
	}
}

type AnalysisResults struct {
	Attributes struct {
		Stats  *AnalysisStats `json:"stats,omitempty"`
		Status string         `json:"status"`
	} `json:"attributes,omitempty"`
	Meta struct {
		FileInfo struct {
			SHA256 string `json:"sha256"`
		} `json:"file_info"`
	} `json:"meta"`
}

type PreviousAnalysisResults struct {
	Attributes struct {
		Stats *AnalysisStats `json:"last_analysis_stats,omitempty"`
	} `json:"attributes,omitempty"`
}

type AnalysisStats struct {
	Suspicious *int `json:"suspicious,omitempty"`
	Malicious  *int `json:"malicious,omitempty"`
}

type ScanResult struct {
	Safe     bool
	Hash     string
	FileName string
}

func ScanFiles(ctx context.Context, files []io.Reader, names []string) ([]ScanResult, error) {
	errs, gctx := errgroup.WithContext(ctx)
	fileCount := len(files)

	c := make(chan ScanResult, fileCount)

	for i := range files {
		count := i
		errs.Go(func() error {
			scanResult, err := scanFile(gctx, files[count], names[count])
			if err != nil {
				return fmt.Errorf("failed to scan %s: %w", names[count], err)
			}
			c <- *scanResult
			return nil
		})
	}

	go func() {
		defer close(c)
		_ = errs.Wait()
	}()

	results := make([]ScanResult, 0)
	for res := range c {
		results = append(results, res)
	}

	if err := errs.Wait(); err != nil {
		return results, fmt.Errorf("failed to scan files: %w", err)
	}

	return results, nil
}

func scanFile(ctx context.Context, file io.Reader, name string) (*ScanResult, error) {
	hash := sha256.New()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed read file: %w", err)
	}

	_, err = hash.Write(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to hash file: %w", err)
	}

	checksum := hash.Sum(nil)
	hashStr := fmt.Sprintf("%x", checksum)

	var previousAnalysisResults PreviousAnalysisResults

	hasPreviousAnalysis := true
	_, err = client.GetData(vt.URL("files/%x", checksum), &previousAnalysisResults)
	if err != nil {
		var vtErr vt.Error
		if errors.As(err, &vtErr) {
			if vtErr.Code != "NotFoundError" {
				return nil, fmt.Errorf("failed to get previous analysis results: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get previous analysis results: %w", err)
		}
		hasPreviousAnalysis = false
	}

	if hasPreviousAnalysis {
		slox.Info(ctx, "file already scanned", slog.String("file", name), slog.String("hash", hashStr))
		if previousAnalysisResults.Attributes.Stats == nil {
			return nil, fmt.Errorf("no stats available on previous analysis")
		}
		safe, err := isResultSafe(*previousAnalysisResults.Attributes.Stats)
		if err != nil {
			return nil, fmt.Errorf("failed to determine if file is safe: %w", err)
		}
		return &ScanResult{
			Safe:     safe,
			Hash:     hashStr,
			FileName: name,
		}, nil
	}

	scan, err := client.NewFileScanner().Scan(bytes.NewReader(fileBytes), name, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}
	analysisID := scan.ID()
	slox.Info(ctx, "uploaded virus scan", slog.String("file", name), slog.String("analysis_id", analysisID))

	var analysisResults AnalysisResults
	for analysisResults.Attributes.Status != "completed" {
		time.Sleep(time.Second * 15)

		_, err := client.GetData(vt.URL("analyses/%s", analysisID), &analysisResults)
		if err != nil {
			return nil, fmt.Errorf("failed to get analysis results: %w", err)
		}
	}

	if analysisResults.Attributes.Stats == nil {
		return nil, fmt.Errorf("no stats available")
	}

	safe, err := isResultSafe(*analysisResults.Attributes.Stats)
	if err != nil {
		return nil, fmt.Errorf("failed to determine if file is safe: %w", err)
	}
	return &ScanResult{
		Safe:     safe,
		Hash:     hashStr,
		FileName: name,
	}, nil
}

func isResultSafe(stats AnalysisStats) (bool, error) {
	if stats.Malicious == nil || stats.Suspicious == nil {
		return false, fmt.Errorf("missing malicious or suspicious stats")
	}

	// Why 1? Well because some company made a shitty AI and it flags random mods.
	return *stats.Malicious <= 1 && *stats.Suspicious <= 1, nil
}
