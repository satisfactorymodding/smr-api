package versionupload

import (
	"fmt"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/validation"
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

	fatalError := func(ctx workflow.Context, err error) error {
		if err != nil {
			cleanCtx, _ := workflow.NewDisconnectedContext(ctx)

			_ = workflow.ExecuteActivity(cleanCtx, VersionUpload.StoreRedisStateActivity, StoreRedisStateArgs{
				UploadID: args.UploadID,
				Err:      err.Error(),
			}).Get(cleanCtx, nil)

			_ = workflow.ExecuteActivity(cleanCtx, VersionUpload.RemoveUploadActivity, RemoveUploadArgs{
				ModID:    args.ModID,
				UploadID: args.UploadID,
			}).Get(cleanCtx, nil)
		}
		return err
	}

	err := workflow.ExecuteActivity(ctx, VersionUpload.CompleteUploadMultipartModActivity, CompleteUploadMultipartModArgs{
		ModID:    args.ModID,
		UploadID: args.UploadID,
	}).Get(ctx, nil)
	if err != nil {
		return fatalError(ctx, err)
	}

	var modInfo validation.ModInfo
	err = workflow.ExecuteActivity(ctx, VersionUpload.ExtractModInfoActivity, ExtractModInfoArgs{
		ModID:    args.ModID,
		UploadID: args.UploadID,
	}).Get(ctx, &modInfo)
	if err != nil {
		return fatalError(ctx, err)
	}

	var dbVersion *ent.Version
	createVersionCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID:            fmt.Sprintf("create-version-%s-%s", modInfo.ModReference, modInfo.Version),
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
	})
	err = workflow.ExecuteChildWorkflow(createVersionCtx, VersionUpload.CreateVersionWorkflow, CreateVersionArgs{
		ModID:          args.ModID,
		UploadID:       args.UploadID,
		Version:        args.Version,
		ModInfo:        modInfo,
		SkipVirusCheck: args.SkipVirusCheck,
	},
	).Get(createVersionCtx, &dbVersion)
	if err != nil {
		if temporal.IsWorkflowExecutionAlreadyStartedError(err) {
			return fatalError(ctx, fmt.Errorf("another upload for %s v%s is already in progress", modInfo.ModReference, modInfo.Version))
		}
		return fatalError(ctx, err)
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
			return err
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

type CreateVersionArgs struct {
	ModID          string               `json:"mod_id"`
	UploadID       string               `json:"upload_id"`
	Version        generated.NewVersion `json:"version"`
	SkipVirusCheck bool                 `json:"skip_virus_check"`
	ModInfo        validation.ModInfo   `json:"mod_info"`
}

func (*A) CreateVersionWorkflow(ctx workflow.Context, args CreateVersionArgs) (*ent.Version, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 10,
		},
	})

	err := workflow.ExecuteActivity(ctx, VersionUpload.ValidateVersionActivity, ValidateVersionArgs{
		ModID:   args.ModID,
		ModInfo: args.ModInfo,
	}).Get(ctx, nil)
	if err != nil {
		// Does not need cleanup, because no changes have been made so far
		return nil, err
	}

	fatalError := func(ctx workflow.Context, err error) error {
		if err != nil {
			cleanCtx, _ := workflow.NewDisconnectedContext(ctx)
			_ = workflow.ExecuteActivity(cleanCtx, VersionUpload.RemoveModActivity, RemoveModArgs{
				ModID:   args.ModID,
				ModInfo: args.ModInfo,
			}).Get(cleanCtx, nil)
		}
		return err
	}

	var metadata *string
	_ = workflow.ExecuteActivity(ctx, VersionUpload.ExtractMetadataActivity, ExtractMetadataArgs{
		ModID:    args.ModID,
		UploadID: args.UploadID,
		ModInfo:  args.ModInfo,
	}).Get(ctx, &metadata)

	var fileKey string
	err = workflow.ExecuteActivity(ctx, VersionUpload.RenameVersionActivity, RenameVersionArgs{
		ModID:    args.ModID,
		UploadID: args.UploadID,
		Version:  args.ModInfo.Version,
	}).Get(ctx, &fileKey)
	if err != nil {
		return nil, fatalError(ctx, err)
	}

	var targetsData []ModTargetData
	err = workflow.ExecuteActivity(ctx, VersionUpload.SeparateModTargetsActivity, SeparateModTargetsArgs{
		ModID:   args.ModID,
		ModInfo: args.ModInfo,
		FileKey: fileKey,
	}).Get(ctx, &targetsData)
	if err != nil {
		return nil, fatalError(ctx, err)
	}

	var dbVersion *ent.Version
	err = workflow.ExecuteActivity(ctx, VersionUpload.CreateVersionInDatabaseActivity, CreateVersionInDatabaseArgs{
		ModID:    args.ModID,
		ModInfo:  args.ModInfo,
		FileKey:  fileKey,
		Targets:  targetsData,
		Version:  args.Version,
		Metadata: metadata,
	}).Get(ctx, &dbVersion)
	if err != nil {
		return nil, fatalError(ctx, err)
	}

	return dbVersion, nil
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
