package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/Vilsol/slox"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/satisfactorymodding/smr-api/ent"
	"github.com/spf13/viper"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/satisfactorymodding/smr-api/ent/runtime"
)

type dbKey struct{}
type txKey struct{}

type dbClient struct {
	Client *ent.Client
}

type Config struct {
	Address string
}

type ShowDB struct{}

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

	realClient := ent.NewClient(
		ent.Driver(poolDriver),
		ent.Log(func(v ...interface{}) {
			fmt.Println(v...)
			slox.Info(ctx, fmt.Sprint(v...))
		}),
	)

	return context.WithValue(ctx, dbKey{}, &dbClient{
		Client: realClient,
	}), nil
}

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
					finalError = errors.Join(finalError, fmt.Errorf("panic when rolling back: %w", err))
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

		return finalError
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed committing transaction: %w", err)
	}

	return nil
}
