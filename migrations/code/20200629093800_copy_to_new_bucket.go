package code

import (
	"context"

	"github.com/lab259/go-migration"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/redis/jobs"
	"github.com/satisfactorymodding/smr-api/storage"
)

func init() {
	migration.NewCodeMigration(
		func(_ interface{}) error {
			ctx, err := db.WithDB(context.Background())
			if err != nil {
				return err
			}
			storage.ScheduleCopyAllObjectsFromOldBucket(func(key string) {
				jobs.SubmitJobCopyObjectFromOldBucketTask(ctx, key)
			})
			return nil
		},
	)
}
