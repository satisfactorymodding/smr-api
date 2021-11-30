package code

import (
	"context"

	"github.com/satisfactorymodding/smr-api/migrations/utils"

	"github.com/lab259/go-migration"
	"github.com/rs/zerolog/log"
)

func init() {
	migration.NewCodeMigration(
		func(executionContext interface{}) error {
			ctx := log.Logger.WithContext(context.TODO())
			utils.ReindexAllModFiles(ctx, true, nil, nil)
			return nil
		},
	)
}
