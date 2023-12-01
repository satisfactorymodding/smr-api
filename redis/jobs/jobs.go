package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Vilsol/slox"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/vmihailenco/taskq/extra/taskqotel/v3"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/redisq"

	"github.com/satisfactorymodding/smr-api/redis/jobs/tasks"
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
		Name:               "api-worker",
		Redis:              connection,
		ReservationTimeout: time.Hour,
	})

	QueueFactory.Range(func(q taskq.Queue) bool {
		consumer := q.Consumer()
		consumer.AddHook(&taskqotel.OpenTelemetryHook{})
		return true
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
		slox.Error(ctx, "error adding task", slog.Any("err", err))
	}
}

func SubmitJobUpdateDBFromModVersionJSONFileTask(ctx context.Context, modID string, version string) {
	task, _ := json.Marshal(tasks.UpdateDBFromModVersionJSONFileData{
		ModID:     modID,
		VersionID: version,
	})

	err := queue.Add(tasks.UpdateDBFromModVersionJSONFileTask.WithArgs(ctx, task))
	if err != nil {
		slox.Error(ctx, "error adding task", slog.Any("err", err))
	}
}

func SubmitJobCopyObjectFromOldBucketTask(ctx context.Context, key string) {
	task, _ := json.Marshal(tasks.CopyObjectFromOldBucketData{
		Key: key,
	})

	err := queue.Add(tasks.CopyObjectFromOldBucketTask.WithArgs(ctx, task))
	if err != nil {
		slox.Error(ctx, "error adding task", slog.Any("err", err))
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
		slox.Error(ctx, "error adding task", slog.Any("err", err))
	}
}

func Purge() {
	err := queue.Purge()
	if err != nil {
		slog.Error("failed purging queue", slog.Any("err", err))
	}
}
