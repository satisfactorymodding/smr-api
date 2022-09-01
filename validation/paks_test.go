package validation

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Vilsol/ue4pak/parser"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestExtractDataFromPak(t *testing.T) {
	paks, err := filepath.Glob("paks/*.pak")

	ctx := context.Background()
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()

	if err != nil {
		panic(err)
	}

	for _, f := range paks {
		fmt.Println("Parsing file:", f)

		data, err := os.ReadFile(f)

		if err != nil {
			panic(err)
		}

		reader := &parser.PakByteReader{
			Bytes: data,
		}

		pakData, err := AttemptExtractDataFromPak(ctx, reader)

		if err != nil {
			log.Err(err).Msg("error parsing pak")
			t.Error(err)
		} else {
			marshal, _ := json.MarshalIndent(pakData, "", "  ")

			if err := os.WriteFile(f+".json", marshal, 0644); err != nil {
				t.Error(err)
			}
		}
	}

	/*
		f, _ := os.Open("E:\\Program Files\\Epic Games\\SatisfactoryExperimental\\FactoryGame\\Content\\Paks\\FactoryGame-WindowsNoEditor.pak")
		pakData, err := AttemptExtractDataFromPak(f)

		if err != nil {
			fmt.Println(err)
		}

		marshal, _ := json.MarshalIndent(pakData, "", "  ")

		ioutil.WriteFile("E:\\Program Files\\Epic Games\\SatisfactoryExperimental\\FactoryGame\\Content\\Paks\\FactoryGame-WindowsNoEditor.pak.json", marshal, 0644)
	*/
}
