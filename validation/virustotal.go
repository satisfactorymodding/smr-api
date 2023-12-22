package validation

import (
	"context"
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
}

func ScanFiles(_ context.Context, files []io.Reader, names []string) (bool, error) {
	errs, gctx := errgroup.WithContext(context.Background())
	fileCount := len(files)

	c := make(chan bool)

	for i := 0; i < fileCount; i++ {
		count := i
		errs.Go(func() error {
			ok, err := scanFile(gctx, files[count], names[count])
			if err != nil {
				return fmt.Errorf("failed to scan file: %w", err)
			}
			c <- ok
			return nil
		})
	}
	go func() {
		_ = errs.Wait()
		close(c)
	}()

	success := true
	for res := range c {
		if !res {
			success = false
			break
		}
	}

	if err := errs.Wait(); err != nil {
		return false, fmt.Errorf("failed to scan file: %w", err)
	}

	return success, nil
}

func scanFile(ctx context.Context, file io.Reader, name string) (bool, error) {
	scan, err := client.NewFileScanner().Scan(file, name, nil)
	if err != nil {
		return false, fmt.Errorf("failed to scan file: %w", err)
	}

	analysisID := scan.ID()

	slox.Info(ctx, "uploaded virus scan", slog.String("file", name), slog.String("analysis_id", analysisID))

	for {
		time.Sleep(time.Second * 15)

		var target AnalysisResults
		_, err = client.GetData(vt.URL("analyses/%s", analysisID), &target)

		if err != nil {
			return false, fmt.Errorf("failed to get analysis results: %w", err)
		}

		if target.Attributes.Status != "completed" {
			continue
		}

		if target.Attributes.Stats == nil {
			slox.Error(ctx, "no stats available", slog.Any("err", err), slog.String("file", name))
			return false, nil
		}

		if target.Attributes.Stats.Malicious == nil || target.Attributes.Stats.Suspicious == nil {
			slox.Error(ctx, "unable to determine malicious or suspicious file", slog.Any("err", err), slog.String("file", name))
			return false, nil
		}

		// Why 1? Well because some company made a shitty AI and it flags random mods.
		if *target.Attributes.Stats.Malicious > 1 || *target.Attributes.Stats.Suspicious > 1 {
			slox.Error(ctx, "suspicious or malicious file found", slog.Any("err", err), slog.String("file", name))
			return false, nil
		}

		break
	}

	return true, nil
}
