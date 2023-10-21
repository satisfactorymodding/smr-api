package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Vilsol/slox"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/satisfactorymodding/smr-api/db/postgres/otel"
)

var (
	db      *gorm.DB
	dbCache *cache.Cache
)

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
	slox.Info(ctx, fmt.Sprintf(msg, data...), slog.String("file", utils.FileWithLineNum()))
}

func (*GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	slox.Warn(ctx, fmt.Sprintf(msg, data...), slog.String("file", utils.FileWithLineNum()))
}

func (*GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	slox.Error(ctx, fmt.Sprintf(msg, data...), slog.String("file", utils.FileWithLineNum()))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	since := time.Since(begin)
	elapsed := float64(since.Nanoseconds()) / 1e6

	sql, rows := fc()

	level := slog.LevelInfo
	attrs := make([]slog.Attr, 0)
	toLog := false
	switch {
	case err != nil:
		level = slog.LevelError
		attrs = append(attrs, slog.Any("err", err))
		toLog = true
	case since > l.SlowThreshold && l.SlowThreshold != 0:
		level = slog.LevelWarn
		toLog = true
	case l.Debug:
		level = slog.LevelInfo
		toLog = true
	}

	if toLog {
		if len(sql) > 256 {
			sql = sql[:256]
		}

		attrs = append(attrs, slog.Float64("elapsed", elapsed))
		attrs = append(attrs, slog.Int64("rows", rows))
		slog.LogAttrs(ctx, level, sql, attrs...)
	}
}

var debugEnabled = false

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
			SlowThreshold: time.Second,
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

	if debugEnabled {
		db = db.Debug()
	}

	dbCache = cache.New(time.Second*5, time.Second*10)

	// TODO Create search indexes

	slox.Info(ctx, "Postgres initialized")
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

func EnableDebug() {
	if db != nil {
		db = db.Debug()
	}

	debugEnabled = true
}
