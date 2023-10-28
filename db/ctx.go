package db

import (
	"context"

	"github.com/Vilsol/slox"
)

// ReWrapCtx re-wraps the old logger and database but with a new context
func ReWrapCtx(ctx context.Context) context.Context {
	newCtx := context.Background()
	newCtx = slox.Into(newCtx, slox.From(ctx))
	newCtx = TransferContext(ctx, newCtx)
	return newCtx
}
