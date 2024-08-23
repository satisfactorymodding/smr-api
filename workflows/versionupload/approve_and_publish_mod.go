package versionupload

import (
	"context"
	"log/slog"
	"time"

	"github.com/Vilsol/slox"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/integrations"
)

type ApproveAndPublishModArgs struct {
	ModID     string `json:"mod_id"`
	VersionID string `json:"version_id"`
}

func (*A) ApproveAndPublishModActivity(ctx context.Context, args ApproveAndPublishModArgs) error {
	slox.Info(ctx, "approving mod", slog.String("mod", args.ModID), slog.String("version", args.VersionID))

	version, err := db.From(ctx).Version.Get(ctx, args.VersionID)
	if err != nil {
		return err
	}

	if err := version.Update().SetApproved(true).Exec(ctx); err != nil {
		return err
	}

	if err := db.From(ctx).Mod.UpdateOneID(args.ModID).SetLastVersionDate(time.Now()).Exec(ctx); err != nil {
		return err
	}

	go integrations.NewVersion(db.ReWrapCtx(ctx), version)

	return nil
}
