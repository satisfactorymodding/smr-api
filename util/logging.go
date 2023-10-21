package util

import (
	"context"

	"github.com/Vilsol/slox"
)

// ReWrapCtx re-wraps the old logger but with a new context
func ReWrapCtx(ctx context.Context) context.Context {
	return slox.Into(context.Background(), slox.From(ctx))
}
