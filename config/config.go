package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Vilsol/slox"
	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
)

var configDir = "."

const (
	ansiReset         = "\033[0m"
	ansiBold          = "\033[1m"
	ansiWhite         = "\033[38m"
	ansiBrightMagenta = "\033[95m"
)

func SetConfigDir(newConfigDir string) {
	configDir = newConfigDir
}

func InitializeConfig(baseCtx context.Context) context.Context {
	viper.SetConfigName("config")
	viper.AddConfigPath(configDir)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("repo")

	initializeDefaults()

	err := viper.ReadInConfig() //nolint:ifshort

	if !viper.GetBool("production") {
		slog.SetDefault(slog.New(
			StackRewriter{
				Upstream: tint.NewHandler(os.Stderr, &tint.Options{
					Level:      slog.LevelDebug,
					AddSource:  true,
					TimeFormat: time.RFC3339Nano,
					ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
						if attr.Key == slog.LevelKey {
							level := attr.Value.Any().(slog.Level)
							if level == slog.LevelDebug {
								attr.Value = slog.StringValue(ansiBrightMagenta + "DBG" + ansiReset)
							}
						} else if attr.Key == slog.MessageKey {
							attr.Value = slog.StringValue(ansiBold + ansiWhite + fmt.Sprint(attr.Value.Any()) + ansiReset)
						}
						return attr
					},
				}).WithAttrs([]slog.Attr{slog.String("service", "api")}),
			},
		))
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		})))
	}

	if baseCtx == nil {
		baseCtx = context.Background()
	}

	if err != nil {
		slog.Warn("config initialized using defaults and environment only!", slog.Any("err", err))
	}

	slox.Info(baseCtx, "Config initialized")

	return baseCtx
}

func initializeDefaults() {
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", "5020")

	viper.SetDefault("production", true)
	viper.SetDefault("profiler", false)

	viper.SetDefault("database.redis.host", "localhost")
	viper.SetDefault("database.redis.port", 6379)
	viper.SetDefault("database.redis.pass", "")
	viper.SetDefault("database.redis.db", 1)
	viper.SetDefault("database.redis.job_db", 2)

	viper.SetDefault("database.postgres.host", "localhost")
	viper.SetDefault("database.postgres.port", 5432)
	viper.SetDefault("database.postgres.user", "postgres")
	viper.SetDefault("database.postgres.pass", "REPLACE_ME")
	viper.SetDefault("database.postgres.db", "postgres")

	viper.SetDefault("storage.type", "s3")
	viper.SetDefault("storage.bucket", "smr")
	viper.SetDefault("storage.key", "REPLACE_ME_KEY")
	viper.SetDefault("storage.secret", "REPLACE_ME_SECRET")
	viper.SetDefault("storage.endpoint", "http://localhost:9000")
	viper.SetDefault("storage.region", "eu-central-1")
	viper.SetDefault("storage.base_url", "http://localhost:9000")
	viper.SetDefault("storage.keypath", "%s/file/%s/%s")

	viper.SetDefault("oauth.github.client_id", "")
	viper.SetDefault("oauth.github.client_secret", "")

	viper.SetDefault("oauth.google.client_id", "")
	viper.SetDefault("oauth.google.client_secret", "")

	viper.SetDefault("oauth.facebook.client_id", "")
	viper.SetDefault("oauth.facebook.client_secret", "")

	viper.SetDefault("paseto.public_key", "")
	viper.SetDefault("paseto.private_key", "")

	viper.SetDefault("discord.webhook_url", "")

	viper.SetDefault("discourse.url", "")
	viper.SetDefault("discourse.sso_secret", "")

	viper.SetDefault("frontend.url", "")

	viper.SetDefault("virustotal.key", "")

	viper.SetDefault("feature_flags.allow_multi_target_upload", false)

	viper.SetDefault("extractor_host", "localhost:50051")
}
