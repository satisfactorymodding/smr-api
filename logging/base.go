package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/viper"
)

const (
	ansiReset         = "\033[0m"
	ansiBold          = "\033[1m"
	ansiWhite         = "\033[38m"
	ansiBrightMagenta = "\033[95m"
)

func SetupLogger(out io.Writer) error {
	var terminalHandler slog.Handler

	if !viper.GetBool("production") {
		if out == nil {
			out = os.Stderr
		}

		terminalHandler = StackRewriter{
			Upstream: tint.NewHandler(out, &tint.Options{
				Level:      slog.LevelDebug,
				AddSource:  true,
				TimeFormat: time.RFC3339Nano,
				ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
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
		}
	} else {
		if out == nil {
			out = os.Stdout
		}

		terminalHandler = StackRewriter{
			Upstream: slog.NewJSONHandler(out, &slog.HandlerOptions{
				AddSource: true,
			}),
		}
	}

	conf := NewClientConfig(map[string]string{
		"app": "ficsit-app-api",
	})

	loki, err := NewClientProto(conf)
	if err != nil {
		return err
	}

	logger := slog.New(slogmulti.Fanout(
		slog.NewJSONHandler(loki, nil),
		terminalHandler,
	))

	slog.SetDefault(logger)

	return nil
}
