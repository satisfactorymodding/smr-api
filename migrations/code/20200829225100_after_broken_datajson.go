package code

import (
	"context"

	"github.com/lab259/go-migration"

	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/migrations/utils"
)

func init() {
	migration.NewCodeMigration(
		func(ctxInt interface{}) error {
			ctx := ctxInt.(context.Context)
			return utils.ReindexAllModFiles(ctx, true, nil, func(version *ent.Version) bool {
				return version.Hash == ""
			})
		},
	)
}
