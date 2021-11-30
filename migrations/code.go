package migrations

import (
	"context"
	"os"

	postgres2 "github.com/satisfactorymodding/smr-api/db/postgres"

	"strings"

	// Import all migrations
	_ "github.com/satisfactorymodding/smr-api/migrations/code"

	"github.com/lab259/go-migration"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type codeMigrationLogger struct {
	log *zerolog.Logger
}

func (c codeMigrationLogger) Write(p []byte) (n int, err error) {
	message := strings.TrimRight(string(p), "\n")
	if len(message) > 0 {
		log.Info().Msg(message)
	}
	return len(p), nil
}

func codeMigrations(ctx context.Context) {
	source := migration.DefaultCodeSource()

	// TODO Custom reporter, this one's very ugly
	reporter := migration.NewDefaultReporterWithParams(codeMigrationLogger{log: log.Ctx(ctx)}, os.Exit)

	db, _ := postgres2.GetDB().DB()
	manager := migration.NewDefaultManager(migration.NewPostgreSQLTarget(db), source)
	runner := migration.NewArgsRunnerCustom(reporter, manager, os.Exit, "migrate")
	runner.Run(db)
}
