package util

import (
	"context"

	"github.com/rs/zerolog"
)

// ReWrapCtx re-wraps the old logger but with a new context
func ReWrapCtx(ctx context.Context) context.Context {
	return zerolog.Ctx(ctx).WithContext(context.Background())
}
