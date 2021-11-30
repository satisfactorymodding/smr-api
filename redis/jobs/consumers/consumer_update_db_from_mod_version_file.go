package consumers

import (
	"context"
	"encoding/json"

	"github.com/satisfactorymodding/smr-api/redis/jobs/tasks"

	"github.com/pkg/errors"

	"github.com/vmihailenco/taskq/v3"
)

func init() {
	tasks.UpdateDBFromModVersionFileTask = taskq.RegisterTask(&taskq.TaskOptions{
		Name:    "consumer_update_db_from_mod_version_file",
		Handler: UpdateDBFromModVersionFileConsumer,
	})
}

func UpdateDBFromModVersionFileConsumer(ctx context.Context, payload []byte) error {
	var task tasks.UpdateDBFromModVersionFileData
	if err := json.Unmarshal(payload, &task); err != nil {
		return errors.Wrap(err, "failed to unmarshal task")
	}
	return UpdateModDataFromStorage(ctx, task.ModID, task.VersionID, true)
}
