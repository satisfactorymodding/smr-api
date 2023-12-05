package validation

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/VirusTotal/vt-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
		Status string `json:"status"`
		Stats  *struct {
			Suspicious *int `json:"suspicious,omitempty"`
			Malicious  *int `json:"malicious,omitempty"`
		} `json:"stats,omitempty"`
	} `json:"attributes,omitempty"`
}

func ScanFiles(ctx context.Context, files []io.Reader, names []string) (bool, error) {
	for i, file := range files {
		scan, err := client.NewFileScanner().Scan(file, names[i], nil)

		if err != nil {
			return false, errors.Wrap(err, "failed to scan file")
		}

		analysisID := scan.ID()

		log.Info().Msgf("uploaded virus scan for file %s and analysis ID: %s", names[i], analysisID)

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
				return false, nil
			}

			if target.Attributes.Stats.Malicious == nil || target.Attributes.Stats.Suspicious == nil {
				return false, nil
			}

			// Why 1? Well because some company made a shitty AI and it flags random mods.
			if *target.Attributes.Stats.Malicious > 1 || *target.Attributes.Stats.Suspicious > 1 {
				log.Error().Msgf("suspicious or malicious file found: %s", name)
				return false, nil
			}

			break
		}
	}

	return true, nil
}
