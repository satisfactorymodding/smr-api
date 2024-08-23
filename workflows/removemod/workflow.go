package removemod

import (
	"time"

	"github.com/satisfactorymodding/smr-api/validation"
	"go.temporal.io/sdk/workflow"
)

type A struct{}

var RemoveMod = &A{}

type RemoveModArgs struct {
	ModID    string              `json:"mod_id"`
	ModInfo  *validation.ModInfo `json:"mod_info"`
	UploadID string              `json:"upload_id"`
}

func (*A) RemoveModWorkflow(ctx workflow.Context, args RemoveModArgs) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	return workflow.ExecuteActivity(ctx, RemoveMod.RemoveModActivity, args).Get(ctx, nil)
}
