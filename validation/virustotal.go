package validation

import (
	"context"
	"io"
	"time"

	"github.com/VirusTotal/vt-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

func ScanFiles(ctx context.Context, files []io.Reader, names []string) (bool, error) {
	errs, gctx := errgroup.WithContext(context.Background())
	fileCount := len(files)

	c := make(chan bool)

	for i := 0; i < fileCount; i++ {
		count := i
		errs.Go(func() error {
			ok, err := scanFile(gctx, files[count], names[count])
			if err != nil {
				return errors.Wrap(err, "failed to scan file")
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
		return false, errors.Wrap(err, "failed to scan file")
	}

	return success, nil
}

func scanFile(ctx context.Context, file io.Reader, name string) (bool, error) {
	scan, err := client.NewFileScanner().Scan(file, name, nil)
	if err != nil {
		return false, errors.Wrap(err, "failed to scan file")
	}

	analysisID := scan.ID()

	log.Info().Msgf("uploaded virus scan for file %s and analysis ID: %s", name, analysisID)

	for {
		time.Sleep(time.Second * 15)

		var target AnalysisResults
		_, err = client.GetData(vt.URL("analyses/%s", analysisID), &target)

		if err != nil {
			return false, errors.Wrap(err, "failed to get analysis results")
		}

		if target.Attributes.Status != "completed" {
			continue
		}

		if target.Attributes.Stats == nil {
			log.Error().Msgf("no stats available. failing file: %s", name)
			return false, nil
		}

		if target.Attributes.Stats.Malicious == nil || target.Attributes.Stats.Suspicious == nil {
			log.Error().Msgf("unable to determine malicious or suspicious File: %s", name)
			return false, nil
		}

		if *target.Attributes.Stats.Malicious > 0 || *target.Attributes.Stats.Suspicious > 0 {
			log.Error().Msgf("suspicious or malicious file found: %s", name)
			return false, nil
		}

		break
	}

	return true, nil
}
