package validation

import (
	"context"
	"crypto/sha256"
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
		Stats *struct {
			Suspicious *int `json:"suspicious,omitempty"`
			Malicious  *int `json:"malicious,omitempty"`
		} `json:"stats,omitempty"`
		Status string `json:"status"`
	} `json:"attributes,omitempty"`
	Meta struct {
		FileInfo struct {
			SHA256 string `json:"sha256"`
		} `json:"file_info"`
	} `json:"meta"`
}

type PreviousAnalysisResults struct {
	Attributes struct {
		Stats *struct {
			Suspicious *int `json:"suspicious,omitempty"`
			Malicious  *int `json:"malicious,omitempty"`
		} `json:"last_analysis_stats,omitempty"`
	} `json:"attributes,omitempty"`
}

type ScanResult struct {
	Safe     bool
	Hash     *string
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
			select {
			case c <- *scanResult:
			case <-gctx.Done():
				return gctx.Err()
			}
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
	scanResult := ScanResult{
		Safe:     false,
		FileName: name,
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("failed to generate hash for file %w", err)
	}

	checksum := hash.Sum(nil)
	var previousAnalysisResults PreviousAnalysisResults

	_, err := client.GetData(vt.URL("files/%x", checksum), &previousAnalysisResults)

	alreadyScanned := false
	analysisID := ""

	if err == nil {
		alreadyScanned = true
		slox.Info(ctx, "file already scanned, skipping upload", slog.String("file", name))
		hash := fmt.Sprintf("%x", checksum)
		scanResult.Hash = &hash
	} else {
		scan, err := client.NewFileScanner().Scan(file, name, nil)
		if err != nil {
			return &scanResult, fmt.Errorf("failed to scan file: %w", err)
		}
		analysisID := scan.ID()
		slox.Info(ctx, "uploaded virus scan", slog.String("file", name), slog.String("analysis_id", analysisID))
	}

	for {
		var analysisResults AnalysisResults
		var malicious int
		var suspicious int

		if !alreadyScanned {
			time.Sleep(time.Second * 15)

			_, err = client.GetData(vt.URL("analyses/%s", analysisID), &analysisResults)
			if err != nil {
				scanResult.Safe = false
				return nil, fmt.Errorf("failed to get analysis results: %w", err)
			}
			scanResult.Hash = &analysisResults.Meta.FileInfo.SHA256

			if !alreadyScanned && analysisResults.Attributes.Status != "completed" {
				continue
			}

			if analysisResults.Attributes.Stats == nil {
				slox.Error(ctx, "no stats available", slog.Any("err", err), slog.String("file", name))
				scanResult.Safe = false
				return &scanResult, nil
			}

			if analysisResults.Attributes.Stats.Malicious == nil || analysisResults.Attributes.Stats.Suspicious == nil {
				slox.Error(ctx, "unable to determine malicious or suspicious file", slog.Any("err", err), slog.String("file", name))
				scanResult.Safe = false
				return &scanResult, nil
			}
			malicious = *analysisResults.Attributes.Stats.Malicious
			suspicious = *analysisResults.Attributes.Stats.Suspicious
		} else {
			malicious = *previousAnalysisResults.Attributes.Stats.Malicious
			suspicious = *previousAnalysisResults.Attributes.Stats.Suspicious
		}

		// Why 1? Well because some company made a shitty AI and it flags random mods.
		if malicious > 1 || suspicious > 1 {
			slox.Error(ctx, "suspicious or malicious file found", slog.Any("err", err), slog.String("file", name))
			scanResult.Safe = false
			return &scanResult, nil
		}

		scanResult.Safe = true
		break
	}

	return &scanResult, nil
}
