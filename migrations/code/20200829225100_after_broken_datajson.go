package code

import (
	"context"

	"github.com/lab259/go-migration"
	"github.com/rs/zerolog/log"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/migrations/utils"
)

func init() {
	migration.NewCodeMigration(
		func(executionContext interface{}) error {
			ctx := log.Logger.WithContext(context.TODO())
			utils.ReindexAllModFiles(ctx, true, nil, func(version postgres.Version) bool {
				return version.Hash == nil || *version.Hash == ""
			})
			return nil
		},
	)
}
