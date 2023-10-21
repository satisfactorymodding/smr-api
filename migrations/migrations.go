package migrations

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Vilsol/slox"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	postgres2 "github.com/satisfactorymodding/smr-api/db/postgres"

	// Import migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(ctx context.Context) {
	databaseMigrations(ctx)
	codeMigrations(ctx)
	slox.Info(ctx, "Migrations Complete")
}

var migrationDir = "./migrations"

func SetMigrationDir(newMigrationDir string) {
	migrationDir = newMigrationDir
}

func databaseMigrations(ctx context.Context) {
	db, _ := postgres2.DBCtx(ctx).DB()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationDir+"/sql", "postgres", driver)
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
	slox.Info(l.ctx, fmt.Sprintf(strings.TrimRight(format, "\n"), v...))
}

func (l SimpleLogger) Verbose() bool {
	return true
}
