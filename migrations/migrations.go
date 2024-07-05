package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/Vilsol/slox"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/viper"

	// Import migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// Import pgx
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed sql/*.sql
var sqlMigrations embed.FS

func RunMigrations(ctx context.Context) {
	databaseMigrations(ctx)
	codeMigrations(ctx)
	slox.Info(ctx, "Migrations Complete")
}

func databaseMigrations(ctx context.Context) {
	connection, err := sql.Open("pgx", fmt.Sprintf(
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

	driver, err := postgres.WithInstance(connection, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	d, err := iofs.New(sqlMigrations, "sql")
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
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
