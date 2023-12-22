package db

import (
	"context"
	"errors"
	"fmt"

	"ariga.io/entcache"
	"github.com/Vilsol/slox"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/smr-api/generated/ent"

	// Required PGX driver
	_ "github.com/jackc/pgx/v5/stdlib"
)

type (
	dbKey struct{}
	txKey struct{}
)

type dbClient struct {
	Client *ent.Client
}

var debugEnabled = false

// WithDB initializes a new database instance and puts it in the context
func WithDB(ctx context.Context) (context.Context, error) {
	slox.Info(ctx, "initializing db")

	poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"sslmode=disable host=%s port=%d user=%s dbname=%s password=%s",
		viper.GetString("database.postgres.host"),
		viper.GetInt("database.postgres.port"),
		viper.GetString("database.postgres.user"),
		viper.GetString("database.postgres.db"),
		viper.GetString("database.postgres.pass"),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to parse database connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	poolDriver := NewPgxPoolDriver(pool)

	cacheDriver := entcache.NewDriver(
		poolDriver,
		entcache.ContextLevel(),
	)

	realClient := ent.NewClient(
		ent.Driver(cacheDriver),
		ent.Log(func(v ...interface{}) {
			slox.Info(ctx, fmt.Sprint(v...))
		}),
	)

	if debugEnabled {
		realClient = realClient.Debug()
	}

	return context.WithValue(ctx, dbKey{}, &dbClient{
		Client: realClient,
	}), nil
}

// From retrieves a database instance from the context
func From(ctx context.Context) *ent.Client {
	tx := ctx.Value(txKey{})
	if tx != nil {
		return tx.(*ent.Tx).Client()
	}

	db := ctx.Value(dbKey{})
	if db == nil {
		return nil
	}
	return db.(*dbClient).Client
}

// TransferContext transfers a database instance from source to target context
func TransferContext(source context.Context, target context.Context) context.Context {
	db := source.Value(dbKey{})
	if db == nil {
		return target
	}
	return context.WithValue(target, dbKey{}, db)
}

func Tx(ctx context.Context, f func(newCtx context.Context, tx *ent.Tx) error, onError func() error) error {
	db := ctx.Value(dbKey{})
	if db == nil {
		return errors.New("db key not found in context")
	}

	client := db.(*dbClient).Client

	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start database transaction: %w", err)
	}

	newCtx := context.WithValue(ctx, txKey{}, tx)

	if err := f(newCtx, tx); err != nil {
		finalError := err

		func() {
			defer func() {
				if err := recover(); err != nil {
					finalError = errors.Join(finalError, fmt.Errorf("panic when rolling back: %w", err.(error)))
				}
			}()

			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				finalError = errors.Join(finalError, fmt.Errorf("failed rolling back transaction: %w", rollbackErr))
			}

			if onError() != nil {
				onErrorErr := onError()
				if onErrorErr != nil {
					finalError = errors.Join(finalError, onErrorErr)
				}
			}
		}()

		return finalError // nolint
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed committing transaction: %w", err)
	}

	return nil
}

func EnableDebug() {
	debugEnabled = true
}
