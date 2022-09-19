package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/satisfactorymodding/smr-api/db/postgres/otel"
)

var db *gorm.DB
var dbCache *cache.Cache

type UserKey struct{}

type GormLogger struct {
	SlowThreshold time.Duration
	Debug         bool
}

func (l *GormLogger) LogMode(mode logger.LogLevel) logger.Interface {
	return &GormLogger{
		SlowThreshold: l.SlowThreshold,
		Debug:         mode >= 4,
	}
}

func (*GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Info().Str("file", utils.FileWithLineNum()).Msgf(msg, data...)
}

func (*GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warn().Str("file", utils.FileWithLineNum()).Msgf(msg, data...)
}

func (*GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Error().Str("file", utils.FileWithLineNum()).Msgf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	since := time.Since(begin)
	elapsed := float64(since.Nanoseconds()) / 1e6

	sql, rows := fc()

	var logEv *zerolog.Event
	switch {
	case err != nil:
		logEv = log.Err(err)
	case since > l.SlowThreshold && l.SlowThreshold != 0:
		logEv = log.Warn()
	case l.Debug:
		logEv = log.Info()
	}

	if logEv != nil {
		logEv.Str("file", utils.FileWithLineNum()).
			Float64("elapsed", elapsed).
			Int64("rows", rows).
			Msg(sql)
	}
}

func InitializePostgres(ctx context.Context) {
	connection := postgres.Open(fmt.Sprintf(
		"sslmode=disable host=%s port=%d user=%s dbname=%s password=%s",
		viper.GetString("database.postgres.host"),
		viper.GetInt("database.postgres.port"),
		viper.GetString("database.postgres.user"),
		viper.GetString("database.postgres.db"),
		viper.GetString("database.postgres.pass"),
	))

	dbInit, err := gorm.Open(connection, &gorm.Config{
		Logger: &GormLogger{
			SlowThreshold: time.Millisecond * 50,
		},
	})
	if err != nil {
		panic(err)
	}

	err = dbInit.Use(otel.NewPlugin())
	if err != nil {
		panic(err)
	}

	db = dbInit

	dbCache = cache.New(time.Second*5, time.Second*10)

	// TODO Create search indexes

	log.Info().Msg("Postgres initialized")
}

func Save(ctx context.Context, object interface{}) {
	DBCtx(ctx).Save(object)
}

func Delete(ctx context.Context, object interface{}) {
	DBCtx(ctx).Delete(object)
	ClearCache()
}

func DeleteForced(ctx context.Context, object interface{}) {
	DBCtx(ctx).Unscoped().Delete(object)
	ClearCache()
}

func DBCtx(ctx context.Context) *gorm.DB {
	if ctx != nil {
		dbCtx := DBFromContext(ctx)
		if dbCtx != nil {
			return dbCtx
		}

		return db.WithContext(ctx)
	}

	return db
}

func ClearCache() {
	dbCache.Flush()
}
