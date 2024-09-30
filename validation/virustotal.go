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

var requestLimit = 4
var dailyRequestLimit = 500
var monthlyRequestLimit = 15500
var analysisURL = "https://www.virustotal.com/gui/file/%s"

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
	URL      *string
	Hash     *string
	FileName string
}

func ScanFiles(ctx context.Context, files []io.Reader, names []string) ([]ScanResult, error) {
	errs, gctx := errgroup.WithContext(ctx)
	fileCount := len(files)

	c := make(chan ScanResult)

	for i := range fileCount {
		count := i
		errs.Go(func() error {
			scanResult, err := scanFile(gctx, files[count], names[count])
			if err != nil {
				return fmt.Errorf("failed to scan %s: %w", names[count], err)
			}
			c <- scanResult
			return nil
		})
	}
	go func() {
		_ = errs.Wait()
		close(c)
	}()
	var results []ScanResult
	for res := range c {
		results = append(results, res)
		// if !res.Safe {
		// 	break
		// }
	}

	if err := errs.Wait(); err != nil {
		scanResult := []ScanResult{
			{
				Safe:     false,
				FileName: "",
			},
		}
		return scanResult, fmt.Errorf("failed to scan files: %w", err)
	}

	return results, nil
}

func scanFile(ctx context.Context, file io.Reader, name string) (ScanResult, error) {
	scanResult := ScanResult{
		Safe:     true,
		FileName: name,
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return scanResult, fmt.Errorf("failed to generate hash for file %w", err)
	}

	checksum := hash.Sum(nil)
	var target PreviousAnalysisResults

	_, err := client.GetData(vt.URL("files/%x", checksum), &target)

	alreadyScanned := false
	analysisID := ""

	if err == nil {
		alreadyScanned = true
		slox.Info(ctx, "file already scanned, skipping upload", slog.String("file", name))
		hash := fmt.Sprintf("%x", checksum)
		scanResult.Hash = &hash
		url := fmt.Sprintf(analysisURL, hash)
		scanResult.URL = &url
	} else {
		scan, err := client.NewFileScanner().Scan(file, name, nil)

		if err != nil {
			return scanResult, fmt.Errorf("failed to scan file: %w", err)
		}
		analysisID := scan.ID()
		slox.Info(ctx, "uploaded virus scan", slog.String("file", name), slog.String("analysis_id", analysisID))
	}

	for {
		if !alreadyScanned {
			var target AnalysisResults

			time.Sleep(time.Second * 15)

			_, err = client.GetData(vt.URL("analyses/%s", analysisID), &target)
			if err != nil {
				scanResult.Safe = false
				return scanResult, fmt.Errorf("failed to get analysis results: %w", err)
			}
			scanResult.Hash = &target.Meta.FileInfo.SHA256
			url := fmt.Sprintf(analysisURL, &target.Meta.FileInfo.SHA256)
			scanResult.URL = &url

			if !alreadyScanned && target.Attributes.Status != "completed" {
				continue
			}

			if target.Attributes.Stats == nil {
				slox.Error(ctx, "no stats available", slog.Any("err", err), slog.String("file", name))
				scanResult.Safe = false
				return scanResult, nil
			}

			if target.Attributes.Stats.Malicious == nil || target.Attributes.Stats.Suspicious == nil {
				slox.Error(ctx, "unable to determine malicious or suspicious file", slog.Any("err", err), slog.String("file", name))
				scanResult.Safe = false
				return scanResult, nil
			}

		}

		// Why 1? Well because some company made a shitty AI and it flags random mods.
		if *target.Attributes.Stats.Malicious > 1 || *target.Attributes.Stats.Suspicious > 1 {
			slox.Error(ctx, "suspicious or malicious file found", slog.Any("err", err), slog.String("file", name))
			scanResult.Safe = false
			return scanResult, nil
		}

		break
	}

	return scanResult, nil
}
