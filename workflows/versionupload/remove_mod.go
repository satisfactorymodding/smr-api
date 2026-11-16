package versionupload

import (
	"context"
	"log/slog"

	"github.com/Vilsol/slox"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/validation"
)

type RemoveUploadArgs struct {
	ModID    string `json:"mod_id"`
	UploadID string `json:"upload_id"`
}

func (*A) RemoveUploadActivity(ctx context.Context, args RemoveUploadArgs) error {
	mod, err := db.From(ctx).Mod.Get(ctx, args.ModID)
	if err != nil {
		slox.Error(ctx, "failed to retrieve mod", slog.Any("err", err))
		return nil
	}

	// TODO: cleanup file parts if failure happened before completing multipart upload

	_ = storage.DeleteMod(ctx, mod.ID, mod.Name, args.UploadID)

	return nil
}

type RemoveModArgs struct {
	ModID   string             `json:"mod_id"`
	ModInfo validation.ModInfo `json:"mod_info"`
}

func (*A) RemoveModActivity(ctx context.Context, args RemoveModArgs) error {
	mod, err := db.From(ctx).Mod.Get(ctx, args.ModID)
	if err != nil {
		slox.Error(ctx, "failed to retrieve mod", slog.Any("err", err))
		return nil
	}

	_ = storage.DeleteMod(ctx, mod.ID, mod.Name, args.ModInfo.Version)
	for _, target := range args.ModInfo.Targets {
		_ = storage.DeleteModTarget(ctx, mod.ID, mod.Name, args.ModInfo.Version, target)
	}

	return nil
}
