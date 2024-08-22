package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/lab259/go-migration"
	"github.com/spf13/viper"

	// Import pgx
	_ "github.com/jackc/pgx/v5/stdlib"
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

	db, err := sql.Open("pgx", fmt.Sprintf(
		"sslmode=disable host=%s port=%d user=%s dbname=%s password=%s",
		viper.GetString("database.postgres.host"),
		viper.GetInt("database.postgres.port"),
		viper.GetString("database.postgres.user"),
		viper.GetString("database.postgres.db"),
		viper.GetString("database.postgres.pass"),
	))
	if err != nil {
		panic(err)
	}

	manager := migration.NewDefaultManager(migration.NewPostgreSQLTarget(db), source)
	runner := migration.NewArgsRunnerCustom(reporter, manager, os.Exit, "migrate")
	runner.Run(ctx)
}
