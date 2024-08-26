package statistics

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/Vilsol/slox"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func (*A) InitializeStatisticsWorkflow(ctx context.Context, c client.Client, taskQueue string) {
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
			Workflow:  Statistics.StatisticsWorkflow,
			TaskQueue: taskQueue,
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
