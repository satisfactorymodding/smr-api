package consumers

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/vmihailenco/taskq/v3"

	"github.com/satisfactorymodding/smr-api/redis/jobs/tasks"
	"github.com/satisfactorymodding/smr-api/storage"
)

func init() {
	tasks.CopyObjectFromOldBucketTask = taskq.RegisterTask(&taskq.TaskOptions{
		Name:    "consumer_copy_object_from_old_bucket",
		Handler: CopyObjectFromOldBucketConsumer,
	})
}

func CopyObjectFromOldBucketConsumer(payload []byte) error {
	var task tasks.CopyObjectFromOldBucketData
	if err := json.Unmarshal(payload, &task); err != nil {
		return errors.Wrap(err, "failed to unmarshal task")
	}
	return storage.CopyObjectFromOldBucket(task.Key)
}
