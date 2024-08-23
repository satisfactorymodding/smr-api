package removemod

import (
	"context"
	"log/slog"

	"github.com/Vilsol/slox"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/storage"
)

func (*A) RemoveModActivity(ctx context.Context, args RemoveModArgs) error {
	mod, err := db.From(ctx).Mod.Get(ctx, args.ModID)
	if err != nil {
		slox.Error(ctx, "failed to retrieve mod", slog.Any("err", err))
		return nil
	}

	// TODO: cleanup file parts if failure happened before completing multipart upload

	_ = storage.DeleteMod(ctx, mod.ID, mod.Name, args.UploadID)
	if args.ModInfo != nil {
		_ = storage.DeleteMod(ctx, mod.ID, mod.Name, args.ModInfo.Version)
		for _, target := range args.ModInfo.Targets {
			_ = storage.DeleteModTarget(ctx, mod.ID, mod.Name, args.ModInfo.Version, target)
		}
	}

	return nil
}
