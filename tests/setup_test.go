package tests

import (
	"context"
	"testing"

	smr "github.com/satisfactorymodding/smr-api/api"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/validation"
)

func TestSetup(_ *testing.T) {
	validation.StaticPath = "../static"
	db.DisableDebug()
	ctx, _ := smr.Initialize(context.Background())
	smr.Migrate(ctx)
	smr.Setup(ctx)
}
