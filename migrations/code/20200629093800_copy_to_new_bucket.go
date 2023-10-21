package code

import (
	"context"

	"github.com/lab259/go-migration"

	"github.com/satisfactorymodding/smr-api/redis/jobs"
	"github.com/satisfactorymodding/smr-api/storage"
)

func init() {
	migration.NewCodeMigration(
		func(executionContext interface{}) error {
			storage.ScheduleCopyAllObjectsFromOldBucket(func(key string) {
				jobs.SubmitJobCopyObjectFromOldBucketTask(context.TODO(), key)
			})
			return nil
		},
	)
}
