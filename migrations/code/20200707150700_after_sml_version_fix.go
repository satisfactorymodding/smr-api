package code

import (
	"context"

	"github.com/lab259/go-migration"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/migrations/utils"
)

func init() {
	migration.NewCodeMigration(
		func(_ interface{}) error {
			ctx, err := db.WithDB(context.Background())
			if err != nil {
				return err
			}
			return utils.ReindexAllModFiles(ctx, false, nil, func(version *ent.Version) bool {
				return version.SmlVersion == ""
			})
		},
	)
}
