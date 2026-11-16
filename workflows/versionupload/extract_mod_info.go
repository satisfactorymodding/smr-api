package versionupload

import (
	"context"
	"log/slog"

	"github.com/Vilsol/slox"
	"go.temporal.io/sdk/temporal"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"
)

type ExtractModInfoArgs struct {
	ModID    string `json:"mod_id"`
	UploadID string `json:"upload_id"`
}

func (*A) ExtractModInfoActivity(ctx context.Context, args ExtractModInfoArgs) (*validation.ModInfo, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, args.ModID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", args.UploadID))

	fileData, err := downloadMod(ctx, mod, args.UploadID)
	if err != nil {
		return nil, err
	}

	modInfo, err := validation.ExtractModInfo(ctx, fileData, true, mod.ModReference)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("failed extracting mod info", "fatal", err)
	}

	if modInfo.ModReference != mod.ModReference {
		return nil, temporal.NewNonRetryableApplicationError(".uplugin mod_reference does not match mod reference", "fatal", nil)
	}

	if modInfo.Type == validation.DataJSON {
		return nil, temporal.NewNonRetryableApplicationError("data.json mods are obsolete and not allowed", "fatal", nil)
	}

	if modInfo.Type == validation.MultiTargetUEPlugin && !util.FlagEnabled(util.FeatureFlagAllowMultiTargetUpload) {
		return nil, temporal.NewNonRetryableApplicationError("multi-target mods are not allowed", "fatal", nil)
	}

	return modInfo, nil
}
