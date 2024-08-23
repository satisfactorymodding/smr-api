package versionupload

import (
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/validation"
	"github.com/satisfactorymodding/smr-api/workflows/removemod"
)

type A struct{}

var VersionUpload = &A{}

type FinalizeVersionUploadArgs struct {
	ModID          string               `json:"mod_id"`
	UploadID       string               `json:"upload_id"`
	Version        generated.NewVersion `json:"version"`
	SkipVirusCheck bool                 `json:"skip_virus_check"`
}

func (*A) FinalizeVersionUploadWorkflow(ctx workflow.Context, args FinalizeVersionUploadArgs) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 10,
		},
	})

	fatalError := func(ctx workflow.Context, err error, modInfo *validation.ModInfo) error {
		if err != nil {
			workflow.ExecuteChildWorkflow(workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
				WorkflowID:        workflow.GetInfo(ctx).WorkflowExecution.ID + "-cleanup",
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			}), removemod.RemoveMod.RemoveModWorkflow, removemod.WorkflowArgs{
				ModID:    args.ModID,
				ModInfo:  modInfo,
				UploadID: args.UploadID,
			}).GetChildWorkflowExecution()

			return workflow.ExecuteActivity(ctx, VersionUpload.StoreRedisStateActivity, StoreRedisStateArgs{
				UploadID: args.UploadID,
				Err:      err.Error(),
			}).Get(ctx, nil)
		}
		return err
	}

	err := workflow.ExecuteActivity(ctx, VersionUpload.CompleteUploadMultipartModActivity, CompleteUploadMultipartModArgs{
		ModID:    args.ModID,
		UploadID: args.UploadID,
	}).Get(ctx, nil)
	if err != nil {
		return fatalError(ctx, err, nil)
	}

	var modInfo validation.ModInfo
	err = workflow.ExecuteActivity(ctx, VersionUpload.ExtractModInfoActivity, ExtractModInfoArgs{
		ModID:    args.ModID,
		UploadID: args.UploadID,
	}).Get(ctx, &modInfo)
	if err != nil {
		return fatalError(ctx, err, nil)
	}

	var metadata *string
	_ = workflow.ExecuteActivity(ctx, VersionUpload.ExtractMetadataActivity, ExtractMetadataArgs{
		ModID:    args.ModID,
		UploadID: args.UploadID,
		ModInfo:  modInfo,
	}).Get(ctx, &metadata)

	var fileKey string
	err = workflow.ExecuteActivity(ctx, VersionUpload.RenameVersionActivity, RenameVersionArgs{
		ModID:    args.ModID,
		UploadID: args.UploadID,
		Version:  modInfo.Version,
	}).Get(ctx, &fileKey)
	if err != nil {
		return fatalError(ctx, err, &modInfo)
	}

	var targetsData []ModTargetData
	err = workflow.ExecuteActivity(ctx, VersionUpload.SeparateModTargetsActivity, SeparateModTargetsArgs{
		ModID:   args.ModID,
		ModInfo: modInfo,
		FileKey: fileKey,
	}).Get(ctx, &targetsData)
	if err != nil {
		return fatalError(ctx, err, &modInfo)
	}

	var dbVersion *ent.Version
	err = workflow.ExecuteActivity(ctx, VersionUpload.CreateVersionInDatabaseActivity, CreateVersionInDatabaseArgs{
		ModID:    args.ModID,
		ModInfo:  modInfo,
		FileKey:  fileKey,
		Targets:  targetsData,
		Version:  args.Version,
		Metadata: metadata,
	}).Get(ctx, &dbVersion)
	if err != nil {
		return fatalError(ctx, err, &modInfo)
	}

	data := &generated.CreateVersionResponse{
		AutoApproved: shouldAutoApprove(modInfo, args.SkipVirusCheck),
		Version:      (*conv.VersionImpl)(nil).Convert(dbVersion),
	}

	if data.AutoApproved {
		err = workflow.ExecuteActivity(ctx, VersionUpload.ApproveAndPublishModActivity, ApproveAndPublishModArgs{
			ModID:     args.ModID,
			VersionID: data.Version.ID,
		}).Get(ctx, nil)
		if err != nil {
			return fatalError(ctx, err, &modInfo)
		}
	}

	err = workflow.ExecuteActivity(ctx, VersionUpload.StoreRedisStateActivity, StoreRedisStateArgs{
		UploadID: args.UploadID,
		Data:     data,
	}).Get(ctx, nil)
	if err != nil {
		return err
	}

	if !data.AutoApproved {
		var scanSuccess bool
		err = workflow.ExecuteActivity(ctx, VersionUpload.ScanModOnVirusTotalActivity, ScanModOnVirusTotalArgs{
			ModID:     args.ModID,
			VersionID: data.Version.ID,
		}).Get(ctx, &scanSuccess)
		if err != nil {
			return err
		}

		if !scanSuccess {
			return nil
		}

		err = workflow.ExecuteActivity(ctx, VersionUpload.ApproveAndPublishModActivity, ApproveAndPublishModArgs{
			ModID:     args.ModID,
			VersionID: data.Version.ID,
		}).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func shouldAutoApprove(modInfo validation.ModInfo, skipVirusCheck bool) bool {
	if skipVirusCheck {
		return true
	}

	for _, obj := range modInfo.Objects {
		if obj.Type != "pak" {
			return false
		}
	}

	return true
}
