package updatemodfromstorage

import (
	"go.temporal.io/sdk/workflow"
)

type A struct{}

var UpdateModFromStorage = &A{}

func (*A) UpdateModDataFromStorageWorkflow(ctx workflow.Context, modID string, versionID string, metadata bool) error {
	return workflow.ExecuteActivity(ctx, UpdateModFromStorage.UpdateModDataFromStorageActivity, modID, versionID, metadata).Get(ctx, nil)
}
