package logging

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
	fs := runtime.CallersFrames([]uintptr{record.PC})
	f, _ := fs.Next()

	valid := func(file string, _ string) bool {
		return !strings.Contains(file, "vilsol/slox") &&
			!strings.Contains(file, "entlogger") &&
			!strings.Contains(file, "/generated/") &&
			!strings.Contains(file, "entgo.io")
	}

	// Early exit
	if valid(f.File, f.Function) {
		//nolint
		return t.Upstream.Handle(ctx, record)
	}

	var pcs [25]uintptr
	runtime.Callers(0, pcs[:])

	start := 0
	for i, pc := range pcs {
		if pc == record.PC {
			start = i
		}
	}

	fs = runtime.CallersFrames(pcs[start:])

	f, _ = fs.Next()
	for f.PC != 0 {
		if !valid(f.File, f.Function) {
			f, _ = fs.Next()
			continue
		}

		record.PC = f.PC
		break
	}

	//nolint
	return t.Upstream.Handle(ctx, record)
}

func (t StackRewriter) WithAttrs(attrs []slog.Attr) slog.Handler {
	return StackRewriter{Upstream: t.Upstream.WithAttrs(attrs)}
}

func (t StackRewriter) WithGroup(name string) slog.Handler {
	return StackRewriter{Upstream: t.Upstream.WithGroup(name)}
}
