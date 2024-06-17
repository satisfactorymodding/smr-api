package code

import (
	"context"
	"strings"

	"github.com/lab259/go-migration"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/ent"

	"github.com/satisfactorymodding/smr-api/migrations/utils"
)

func init() {
	migration.NewCodeMigration(
		func(executionContext interface{}) error {
			ctx, err := db.WithDB(context.Background())
			if err != nil {
				return err
			}
			return utils.ReindexAllModFiles(ctx, true, nil, func(version *ent.Version) bool {
				smlVersion := version.SmlVersion
				return strings.Contains(smlVersion, "3.6.1") || strings.Contains(smlVersion, "3.7.0")
			})
		},
	)
}
