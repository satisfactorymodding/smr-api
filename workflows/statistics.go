package workflows

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"time"

	"github.com/Vilsol/slox"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/redis"
)

var keyRegex = regexp.MustCompile(`^([^:]+):([^:]+):([^:]+):([^:]+)$`)

func initializeStatisticsWorkflow(ctx context.Context, c client.Client) {
	if !viper.GetBool("statistics.enabled") {
		return
	}

	scheduleHandle, err := c.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID: "statistics_update_minutely",
		Spec: client.ScheduleSpec{
			Intervals: []client.ScheduleIntervalSpec{
				{
					Every: time.Minute,
				},
			},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        "statistics_update",
			Workflow:  statisticsWorkflow,
			TaskQueue: RepoTaskQueue,
		},
	})
	if err != nil {
		if !errors.Is(err, temporal.ErrScheduleAlreadyRunning) {
			slox.Error(ctx, "unable to create statistics schedule", slog.Any("err", err))
			os.Exit(1)
		}

		scheduleHandle = c.ScheduleClient().GetHandle(ctx, "statistics_update_minutely")
	}

	if _, err := scheduleHandle.Describe(ctx); err != nil {
		slox.Error(ctx, "unable to register statistics schedule", slog.Any("err", err))
		os.Exit(1)
	}
}

func statisticsWorkflow(ctx workflow.Context) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	return workflow.ExecuteActivity(ctx, updateStatisticsActivity).Get(ctx, nil)
}

func updateStatisticsActivity(ctx context.Context) error {
	start := time.Now()
	keys := redis.GetAllKeys()
	slox.Info(ctx, "statistics fetched", slog.Int("keys", len(keys)), slog.Duration("took", time.Since(start)))
	resultMap := make(map[string]map[string]map[string]uint)
	for _, key := range keys {
		if matches := keyRegex.FindStringSubmatch(key); matches != nil {
			entityType := matches[1]
			entityID := matches[2]
			action := matches[3]

			if _, ok := resultMap[entityType]; !ok {
				resultMap[entityType] = make(map[string]map[string]uint)
			}

			if _, ok := resultMap[entityType][action]; !ok {
				resultMap[entityType][action] = make(map[string]uint)
			}

			resultMap[entityType][action][entityID]++
		}
	}

	for entityType, entityValue := range resultMap {
		for action, actionValue := range entityValue {
			for entityID, count := range actionValue {
				err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
					switch entityType {
					case "mod":
						if action == "view" {
							mod, err := tx.Mod.Get(ctx, entityID)
							if err != nil {
								return err
							}

							if mod != nil {
								currentHotness := mod.Hotness
								if currentHotness > 4 {
									// Preserve some of the hotness
									currentHotness /= 4
								}

								return mod.Update().SetHotness(currentHotness + count).Exec(ctx)
							}
						}
					case "version":
						if action == "download" {
							version, err := tx.Version.Get(ctx, entityID)
							if err != nil {
								return err
							}

							if version != nil {
								currentHotness := version.Hotness
								if currentHotness > 4 {
									// Preserve some of the popularity
									currentHotness /= 4
								}
								return version.Update().SetHotness(currentHotness + count).Exec(ctx)
							}
						}
					}

					return nil
				}, nil)
				if err != nil {
					slox.Error(ctx, "failed updating statistics", slog.Any("err", err))
				}
			}
		}
	}

	type Result struct {
		ModID     string `json:"mod_id"`
		Hotness   uint   `json:"hotness"`
		Downloads uint   `json:"downloads"`
	}

	var resultRows []Result

	err := db.From(ctx).Version.
		Query().
		GroupBy("mod_id").
		Aggregate(
			ent.As(ent.Sum(version.FieldHotness), "hotness"),
			ent.As(ent.Sum(version.FieldDownloads), "downloads"),
		).
		Scan(ctx, &resultRows)
	if err != nil {
		return fmt.Errorf("failed summing version data: %w", err)
	}

	for _, row := range resultRows {
		err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
			mod, err := tx.Mod.Get(ctx, row.ModID)

			if mod == nil || mod.ID == "" {
				return nil
			}

			if err != nil {
				return err
			}

			if mod != nil {
				currentPopularity := mod.Popularity
				if currentPopularity > 4 {
					// Preserve some of the popularity
					currentPopularity /= 4
				}
				return mod.Update().SetPopularity(currentPopularity + row.Hotness).SetDownloads(row.Downloads).Exec(ctx)
			}

			return nil
		}, nil)
		if err != nil {
			slox.Error(ctx, "failed updating mod data", slog.Any("err", err))
			continue
		}
	}

	slox.Info(ctx, "statistics updated", slog.Duration("took", time.Since(start)))

	return nil
}
