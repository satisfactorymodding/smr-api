package config

import (
	"context"
	"log/slog"
	"runtime"
	"strings"
)

var _ slog.Handler = (*StackRewriter)(nil)

type StackRewriter struct {
	Upstream slog.Handler
}

func (t StackRewriter) Enabled(ctx context.Context, level slog.Level) bool {
	return t.Upstream.Enabled(ctx, level)
}

func (t StackRewriter) Handle(ctx context.Context, record slog.Record) error {
	if record.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, _ := fs.Next()
		if strings.HasPrefix(f.Function, "github.com/Vilsol/slox") {
			// skip one more
			var pcs [1]uintptr
			runtime.Callers(5, pcs[:])
			record.PC = pcs[0]
		}
	}

	//nolint
	return t.Upstream.Handle(ctx, record)
}

func (t StackRewriter) WithAttrs(attrs []slog.Attr) slog.Handler {
	return t.Upstream.WithAttrs(attrs)
}

func (t StackRewriter) WithGroup(name string) slog.Handler {
	return t.Upstream.WithGroup(name)
}
