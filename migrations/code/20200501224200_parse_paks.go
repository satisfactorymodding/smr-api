package code

import (
	"context"

	"github.com/lab259/go-migration"

	"github.com/satisfactorymodding/smr-api/migrations/utils"
)

func init() {
	migration.NewCodeMigration(
		func(executionContext interface{}) error {
			utils.ReindexAllModFiles(context.TODO(), true, nil, nil)
			return nil
		},
	)
}
