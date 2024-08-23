package versionupload

import (
	"context"
	"log/slog"

	"github.com/Vilsol/slox"
	"go.temporal.io/sdk/temporal"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/storage"
)

type RenameVersionArgs struct {
	ModID    string `json:"mod_id"`
	UploadID string `json:"upload_id"`
	Version  string `json:"version"`
}

func (*A) RenameVersionActivity(ctx context.Context, args RenameVersionArgs) (string, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, args.ModID)
	if err != nil {
		return "", err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", args.UploadID), slog.String("version", args.Version))

	key, err := storage.RenameVersion(ctx, mod.ID, mod.Name, args.UploadID, args.Version)
	if err != nil {
		return "", temporal.NewNonRetryableApplicationError("failed to upload mod", "fatal", err)
	}
	return key, nil
}
