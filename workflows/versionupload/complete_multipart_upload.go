package versionupload

import (
	"context"
	"log/slog"

	"github.com/Vilsol/slox"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/storage"
)

type CompleteUploadMultipartModArgs struct {
	ModID    string `json:"mod_id"`
	UploadID string `json:"upload_id"`
}

func (*A) CompleteUploadMultipartModActivity(ctx context.Context, args CompleteUploadMultipartModArgs) error {
	mod, err := db.From(ctx).Mod.Get(ctx, args.ModID)
	if err != nil {
		return err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", args.UploadID))

	slox.Info(ctx, "Completing multipart upload")
	_, err = storage.CompleteUploadMultipartMod(ctx, mod.ID, mod.Name, args.UploadID)
	return err
}
