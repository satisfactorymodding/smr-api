package migrations

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/lab259/go-migration"

	postgres2 "github.com/satisfactorymodding/smr-api/db/postgres"

	// Import all migrations
	_ "github.com/satisfactorymodding/smr-api/migrations/code"
)

type codeMigrationLogger struct{}

func (c codeMigrationLogger) Write(p []byte) (int, error) {
	message := strings.TrimRight(string(p), "\n")
	if len(message) > 0 {
		slog.Info(message)
	}
	return len(p), nil
}

func codeMigrations(ctx context.Context) {
	source := migration.DefaultCodeSource()

	// TODO Custom reporter, this one's very ugly
	reporter := migration.NewDefaultReporterWithParams(codeMigrationLogger{}, os.Exit)

	db, _ := postgres2.DBCtx(ctx).DB()
	manager := migration.NewDefaultManager(migration.NewPostgreSQLTarget(db), source)
	runner := migration.NewArgsRunnerCustom(reporter, manager, os.Exit, "migrate")
	runner.Run(db)
}
