package code

import (
	"context"
	"strings"

	"github.com/lab259/go-migration"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiondependency"
	"github.com/satisfactorymodding/smr-api/migrations/utils"
)

func init() {
	migration.NewCodeMigration(
		func(ctxInt interface{}) error {
			ctx := ctxInt.(context.Context)
			return utils.ReindexAllModFiles(ctx, true, nil, func(version *ent.Version) bool {
				smlDependency, err := db.From(ctx).VersionDependency.Query().Where(
					versiondependency.VersionID(version.ID),
					versiondependency.ModID("SML"),
				).First(ctx)
				if err != nil {
					return false
				}
				if smlDependency == nil {
					return false
				}
				return strings.Contains(smlDependency.Condition, "3.6.1") || strings.Contains(smlDependency.Condition, "3.7.0")
			})
		},
	)
}
