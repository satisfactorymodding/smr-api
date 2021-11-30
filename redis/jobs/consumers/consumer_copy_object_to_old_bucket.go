package consumers

import (
	"encoding/json"

	"github.com/satisfactorymodding/smr-api/redis/jobs/tasks"
	"github.com/satisfactorymodding/smr-api/storage"

	"github.com/pkg/errors"

	"github.com/vmihailenco/taskq/v3"
)

func init() {
	tasks.CopyObjectToOldBucketTask = taskq.RegisterTask(&taskq.TaskOptions{
		Name:    "consumer_copy_object_to_old_bucket",
		Handler: CopyObjectToOldBucketConsumer,
	})
}

func CopyObjectToOldBucketConsumer(payload []byte) error {
	var task tasks.CopyObjectToOldBucketData
	if err := json.Unmarshal(payload, &task); err != nil {
		return errors.Wrap(err, "failed to unmarshal task")
	}
	return storage.CopyObjectToOldBucket(task.Key)
}
