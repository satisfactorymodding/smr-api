package tests

import (
	"context"
	smr "github.com/satisfactorymodding/smr-api/api"
	"github.com/satisfactorymodding/smr-api/validation"
	"testing"
)

func TestSetup(t *testing.T) {
	validation.StaticPath = "../static"
	ctx, _ := smr.Initialize(context.Background())
	smr.Migrate(ctx)
	smr.Setup(ctx)
}
