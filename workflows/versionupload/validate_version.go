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
	"github.com/satisfactorymodding/smr-api/db/schema"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	version2 "github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"
)

type ValidateVersionArgs struct {
	ModID   string             `json:"mod_id"`
	ModInfo validation.ModInfo `json:"mod_info"`
}

func (*A) ValidateVersionActivity(ctx context.Context, args ValidateVersionArgs) error {
	mod, err := db.From(ctx).Mod.Get(ctx, args.ModID)
	if err != nil {
		return err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version", args.ModInfo.Version))

	existingPotentiallySoftDeletedVersion, err := db.From(ctx).Version.Query().
		Where(
			version2.ModID(mod.ID),
			version2.Version(args.ModInfo.Version),
		).Only(schema.SkipSoftDelete(ctx))
	if err != nil {
		if !ent.IsNotFound(err) {
			return temporal.NewNonRetryableApplicationError("database error", "fatal", err)
		}
		// 'Not Found' is an okay case (no version with the same version number)
	} else {
		if existingPotentiallySoftDeletedVersion.DeletedAt.IsZero() {
			return temporal.NewNonRetryableApplicationError("this mod already has a published version with this name", "fatal", nil)
		}
		return temporal.NewNonRetryableApplicationError("reusing the version name of a deleted mod version is not allowed", "fatal", nil)
	}

	versions, err := db.From(ctx).Version.Query().
		Order(version2.ByCreatedAt(sql.OrderAsc())).
		Where(version2.ModID(mod.ID), version2.CreatedAtGTE(time.Now().Add(time.Hour*24*-1))).
		All(ctx)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("database error", "fatal", err)
	}

	if len(versions) >= util.VersionsPerModPer24h {
		timeToWait := time.Until(versions[0].CreatedAt.Add(time.Hour * 24)).Minutes()
		return temporal.NewNonRetryableApplicationError(fmt.Sprintf("please wait %.0f minutes to post another version", timeToWait), "fatal", nil)
	}

	return nil
}
