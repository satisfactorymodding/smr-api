package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/satisfactorymodding/smr-api/redis/jobs/tasks"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/redisq"
)

var queue taskq.Queue

func InitializeJobs(ctx context.Context) {
	// TODO Somehow add the logger to taskq

	connection := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("database.redis.host") + ":" + fmt.Sprint(viper.GetInt("database.redis.port")),
		Password: viper.GetString("database.redis.pass"),
		DB:       viper.GetInt("database.redis.job_db"),
	})

	QueueFactory := redisq.NewFactory()

	queue = QueueFactory.RegisterQueue(&taskq.QueueOptions{
		Name:  "api-worker",
		Redis: connection,
	})

	if err := QueueFactory.StartConsumers(ctx); err != nil {
		panic(err)
	}
}

func SubmitJobUpdateDBFromModVersionFileTask(ctx context.Context, modID string, version string) {
	task, _ := json.Marshal(tasks.UpdateDBFromModVersionFileData{
		ModID:     modID,
		VersionID: version,
	})

	err := queue.Add(tasks.UpdateDBFromModVersionFileTask.WithArgs(ctx, task))

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error adding task")
	}
}

func SubmitJobUpdateDBFromModVersionJSONFileTask(ctx context.Context, modID string, version string) {
	task, _ := json.Marshal(tasks.UpdateDBFromModVersionJSONFileData{
		ModID:     modID,
		VersionID: version,
	})

	err := queue.Add(tasks.UpdateDBFromModVersionJSONFileTask.WithArgs(ctx, task))

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error adding task")
	}
}

func SubmitJobCopyObjectFromOldBucketTask(ctx context.Context, key string) {
	task, _ := json.Marshal(tasks.CopyObjectFromOldBucketData{
		Key: key,
	})

	err := queue.Add(tasks.CopyObjectFromOldBucketTask.WithArgs(ctx, task))

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error adding task")
	}
}

func SubmitJobCopyObjectToOldBucketTask(ctx context.Context, key string) {
	task, _ := json.Marshal(tasks.CopyObjectToOldBucketData{
		Key: key,
	})

	err := queue.Add(tasks.CopyObjectToOldBucketTask.WithArgs(ctx, task))

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error adding task")
	}
}

func SubmitJobScanModOnVirusTotalTask(ctx context.Context, modID string, version string, approveAfter bool) {
	task, _ := json.Marshal(tasks.ScanModOnVirusTotalData{
		ModID:        modID,
		VersionID:    version,
		ApproveAfter: approveAfter,
	})

	err := queue.Add(tasks.ScanModOnVirusTotalTask.WithArgs(ctx, task))

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error adding task")
	}
}
