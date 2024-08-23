package removemod

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/satisfactorymodding/smr-api/validation"
)

type A struct{}

var RemoveMod = &A{}

type WorkflowArgs struct {
	ModID    string              `json:"mod_id"`
	ModInfo  *validation.ModInfo `json:"mod_info"`
	UploadID string              `json:"upload_id"`
}

func (*A) RemoveModWorkflow(ctx workflow.Context, args WorkflowArgs) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	return workflow.ExecuteActivity(ctx, RemoveMod.RemoveModActivity, args).Get(ctx, nil)
}
