package migrations

import (
	"context"
	"errors"
	"strings"

	postgres2 "github.com/satisfactorymodding/smr-api/db/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// Import migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

func RunMigrations(ctx context.Context) {
	databaseMigrations(ctx)
	codeMigrations(ctx)
	log.Info().Msg("Migrations Complete")
}

func databaseMigrations(ctx context.Context) {
	db, _ := postgres2.DBCtx(ctx).DB()
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations/sql", "postgres", driver)

	if err != nil {
		panic(err)
	}

	m.Log = &SimpleLogger{
		ctx: ctx,
	}

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}

type SimpleLogger struct {
	ctx context.Context
}

func (l SimpleLogger) Printf(format string, v ...interface{}) {
	log.Ctx(l.ctx).Info().Msgf(strings.TrimRight(format, "\n"), v...)
}

func (l SimpleLogger) Verbose() bool {
	return true
}
