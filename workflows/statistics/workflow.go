package statistics

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type A struct{}

var Statistics = &A{}

func (*A) StatisticsWorkflow(ctx workflow.Context) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	return workflow.ExecuteActivity(ctx, Statistics.UpdateStatisticsActivity).Get(ctx, nil)
}
