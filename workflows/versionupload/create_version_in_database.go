package versionupload

import (
	"context"
	"log/slog"

	"github.com/Vilsol/slox"
	"go.temporal.io/sdk/temporal"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	mod2 "github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"
)

type CreateVersionInDatabaseArgs struct {
	ModID    string               `json:"mod_id"`
	ModInfo  validation.ModInfo   `json:"mod_info"`
	FileKey  string               `json:"file_key"`
	Targets  []ModTargetData      `json:"targets"`
	Version  generated.NewVersion `json:"version"`
	Metadata *string              `json:"metadata"`
}

func (*A) CreateVersionInDatabaseActivity(ctx context.Context, args CreateVersionInDatabaseArgs) (*ent.Version, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, args.ModID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version", args.ModInfo.Version))

	var dbVersion *ent.Version
	if err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		dbVersion, err = tx.Version.Create().
			SetVersion(args.ModInfo.Version).
			SetGameVersion(args.ModInfo.GameVersion).
			SetRequiredOnRemote(args.ModInfo.RequiredOnRemote).
			SetChangelog(args.Version.Changelog).
			SetModID(args.ModID).
			SetStability(util.Stability(args.Version.Stability)).
			SetModReference(args.ModInfo.ModReference).
			SetKey(args.FileKey).
			SetSize(args.ModInfo.Size).
			SetHash(args.ModInfo.Hash).
			SetVersionMajor(int(args.ModInfo.Semver.Major())).
			SetVersionMinor(int(args.ModInfo.Semver.Minor())).
			SetVersionPatch(int(args.ModInfo.Semver.Patch())).
			SetNillableMetadata(args.Metadata).
			Save(ctx)
		if err != nil {
			return err
		}

		for modReference, condition := range args.ModInfo.Dependencies {
			modDependency, err := tx.Mod.Query().Where(mod2.ModReference(modReference)).First(ctx)
			if err != nil {
				return err
			}

			_, err = tx.VersionDependency.Create().
				SetVersion(dbVersion).
				SetModID(modDependency.ID).
				SetCondition(condition).
				SetOptional(false).
				Save(ctx)
			if err != nil {
				return err
			}
		}

		for modReference, condition := range args.ModInfo.OptionalDependencies {
			modDependency, err := tx.Mod.Query().Where(mod2.ModReference(modReference)).First(ctx)
			if err != nil {
				return err
			}

			_, err = tx.VersionDependency.Create().
				SetVersion(dbVersion).
				SetModID(modDependency.ID).
				SetCondition(condition).
				SetOptional(true).
				Save(ctx)
			if err != nil {
				return err
			}
		}

		for _, target := range args.Targets {
			_, err = tx.VersionTarget.Create().
				SetVersion(dbVersion).
				SetTargetName(target.TargetName).
				SetKey(target.Key).
				SetHash(target.Hash).
				SetSize(target.Size).
				Save(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	}, nil); err != nil {
		return nil, temporal.NewNonRetryableApplicationError("database error", "fatal", err)
	}

	return dbVersion, nil
}
