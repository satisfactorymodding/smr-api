package versionupload

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Vilsol/slox"
	"go.temporal.io/sdk/temporal"

	"github.com/satisfactorymodding/smr-api/db"
	version2 "github.com/satisfactorymodding/smr-api/generated/ent/version"
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

	count, err := db.From(ctx).Version.Query().
		Where(version2.ModID(mod.ID), version2.Version(modInfo.Version)).
		Count(ctx)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("database error", "fatal", err)
	}

	if count > 0 {
		return nil, temporal.NewNonRetryableApplicationError("this mod already has a version with this name", "fatal", nil)
	}

	// Allow only new 5 versions per 24h
	versions, err := db.From(ctx).Version.Query().
		Order(version2.ByCreatedAt(sql.OrderAsc())).
		Where(version2.ModID(mod.ID), version2.CreatedAt(time.Now().Add(time.Hour*24*-1))).
		All(ctx)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("database error", "fatal", err)
	}

	if len(versions) >= 5 {
		timeToWait := time.Until(versions[0].CreatedAt.Add(time.Hour * 24)).Minutes()
		return nil, temporal.NewNonRetryableApplicationError(fmt.Sprintf("please wait %.0f minutes to post another version", timeToWait), "fatal", nil)
	}

	return modInfo, nil
}
