package gql

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Vilsol/slox"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/db/schema"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	version2 "github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiondependency"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiontarget"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/redis/jobs"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"
)

func FinalizeVersionUploadAsync(ctx context.Context, mod *ent.Mod, versionID string, version generated.NewVersion) (*generated.CreateVersionResponse, error) {
	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version_id", versionID))

	slox.Info(ctx, "Completing multipart upload")
	success, _ := storage.CompleteUploadMultipartMod(ctx, mod.ID, mod.Name, versionID)

	if !success {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		slox.Error(ctx, "failed uploading mod")
		return nil, errors.New("failed uploading mod")
	}

	modFile, err := storage.GetMod(mod.ID, mod.Name, versionID)
	if err != nil {
		slox.Error(ctx, "failed getting mod", slog.Any("err", err))
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, err
	}

	// TODO Optimize
	fileData, err := io.ReadAll(modFile)
	if err != nil {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		slox.Error(ctx, "failed reading mod file", slog.Any("err", err))
		return nil, fmt.Errorf("failed reading mod file: %w", err)
	}

	modInfo, err := validation.ExtractModInfo(ctx, fileData, true, true, mod.ModReference)
	if err != nil {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		slox.Error(ctx, "failed extracting mod info", slog.Any("err", err))
		return nil, fmt.Errorf("failed extracting mod info: %w", err)
	}

	if modInfo.ModReference != mod.ModReference {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		slox.Error(ctx, "data.json mod_reference does not match mod reference", slog.Any("err", err))
		return nil, errors.New("data.json mod_reference does not match mod reference")
	}

	if modInfo.Type == validation.DataJSON {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		slox.Error(ctx, "data.json mods are obsolete and not allowed", slog.Any("err", err))
		return nil, errors.New("data.json mods are obsolete and not allowed")
	}

	if modInfo.Type == validation.MultiTargetUEPlugin && !util.FlagEnabled(util.FeatureFlagAllowMultiTargetUpload) {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		slox.Error(ctx, "multi-target mods are not allowed", slog.Any("err", err))
		return nil, errors.New("multi-target mods are not allowed")
	}

	versionMajor := int(modInfo.Semver.Major())
	versionMinor := int(modInfo.Semver.Minor())
	versionPatch := int(modInfo.Semver.Patch())

	autoApproved := true
	for _, obj := range modInfo.Objects {
		if obj.Type != "pak" {
			autoApproved = false
			break
		}
	}

	autoApproved = autoApproved || viper.GetBool("skip-virus-check")

	count, err := db.From(ctx).Version.Query().
		Where(version2.ModID(mod.ID), version2.Version(modInfo.Version)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("this mod already has a version with this name")
	}

	// Allow only new 5 versions per 24h
	versions, err := db.From(ctx).Version.Query().
		Order(version2.ByCreatedAt(sql.OrderAsc())).
		Where(version2.ModID(mod.ID), version2.CreatedAt(time.Now().Add(time.Hour*24*-1))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	if len(versions) >= 5 {
		timeToWait := time.Until(versions[0].CreatedAt.Add(time.Hour * 24)).Minutes()
		return nil, fmt.Errorf("please wait %.0f minutes to post another version", timeToWait)
	}

	dbVersion, err := db.From(ctx).Version.Create().
		SetVersion(modInfo.Version).
		SetSmlVersion(modInfo.SMLVersion).
		SetChangelog(version.Changelog).
		SetModID(mod.ID).
		SetStability(util.Stability(version.Stability)).
		SetModReference(modInfo.ModReference).
		SetSize(modInfo.Size).
		SetHash(modInfo.Hash).
		SetVersionMajor(versionMajor).
		SetVersionMinor(versionMinor).
		SetVersionPatch(versionPatch).
		SetApproved(autoApproved).
		Save(ctx)
	if err != nil {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, err
	}

	for modID, condition := range modInfo.Dependencies {
		_, err = db.From(ctx).VersionDependency.Create().
			SetVersion(dbVersion).
			SetModID(modID).
			SetCondition(condition).
			SetOptional(false).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	for modID, condition := range modInfo.OptionalDependencies {
		_, err = db.From(ctx).VersionDependency.Create().
			SetVersion(dbVersion).
			SetModID(modID).
			SetCondition(condition).
			SetOptional(true).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	jsonData, err := json.Marshal(modInfo.Metadata)
	if err != nil {
		slox.Error(ctx, "failed serializing", slog.Any("err", err), slog.String("version_id", dbVersion.ID))
	} else {
		metadata := string(jsonData)
		if _, err := dbVersion.Update().SetMetadata(metadata).Save(ctx); err != nil {
			return nil, err
		}
	}

	if modInfo.Type == validation.MultiTargetUEPlugin {
		targets := make([]*ent.VersionTarget, 0)

		for _, target := range modInfo.Targets {
			dbVersionTarget, err := db.From(ctx).VersionTarget.Create().
				SetVersionID(dbVersion.ID).
				SetTargetName(target).
				Save(ctx)
			if err != nil {
				return nil, err
			}

			targets = append(targets, dbVersionTarget)
		}

		separateSuccess := true
		for _, target := range targets {
			slox.Info(ctx, "separating mod", slog.String("target", target.TargetName), slog.String("mod", mod.Name), slog.String("version", dbVersion.Version))
			success, key, hash, size := storage.SeparateModTarget(ctx, fileData, mod.ID, mod.Name, dbVersion.Version, target.TargetName)

			if !success {
				separateSuccess = false
				break
			}

			if _, err := target.Update().SetKey(key).SetHash(hash).SetSize(size).Save(ctx); err != nil {
				return nil, err
			}
		}

		if !separateSuccess {
			removeMod(ctx, modInfo, mod, dbVersion)

			slox.Error(ctx, "failed to separate mod")
			return nil, errors.New("failed to separate mod")
		}
	}

	success, key := storage.RenameVersion(ctx, mod.ID, mod.Name, versionID, modInfo.Version)

	if !success {
		removeMod(ctx, modInfo, mod, dbVersion)

		slox.Error(ctx, "failed to upload mod")
		return nil, errors.New("failed to upload mod")
	}

	if modInfo.Type == validation.UEPlugin {
		if _, err := db.From(ctx).VersionTarget.Create().
			SetVersionID(dbVersion.ID).
			SetTargetName("Windows").
			SetKey(key).
			SetHash(dbVersion.Hash).
			SetSize(dbVersion.Size).
			Save(ctx); err != nil {
			return nil, err
		}
	}

	if _, err := dbVersion.Update().SetKey(key).Save(ctx); err != nil {
		return nil, err
	}

	if autoApproved {
		if _, err := mod.Update().SetLastVersionDate(time.Now()).Save(ctx); err != nil {
			return nil, err
		}

		go integrations.NewVersion(db.ReWrapCtx(ctx), dbVersion)
	} else {
		slox.Info(ctx, "Submitting version job for virus scan")
		jobs.SubmitJobScanModOnVirusTotalTask(ctx, mod.ID, dbVersion.ID, true)
	}

	return &generated.CreateVersionResponse{
		AutoApproved: autoApproved,
		Version:      (*conv.VersionImpl)(nil).Convert(dbVersion),
	}, nil
}

func removeMod(ctx context.Context, modInfo *validation.ModInfo, mod *ent.Mod, dbVersion *ent.Version) {
	for modID, condition := range modInfo.Dependencies {
		if _, err := db.From(ctx).VersionDependency.Delete().Where(
			versiondependency.VersionID(dbVersion.ID),
			versiondependency.ModID(modID),
			versiondependency.Condition(condition),
			versiondependency.Optional(false),
		).Exec(schema.SkipSoftDelete(ctx)); err != nil {
			slox.Error(ctx, "failed deleting version dependency", slog.Any("err", err))
			return
		}
	}

	for modID, condition := range modInfo.OptionalDependencies {
		if _, err := db.From(ctx).VersionDependency.Delete().Where(
			versiondependency.VersionID(dbVersion.ID),
			versiondependency.ModID(modID),
			versiondependency.Condition(condition),
			versiondependency.Optional(true),
		).Exec(schema.SkipSoftDelete(ctx)); err != nil {
			slox.Error(ctx, "failed deleting version dependency", slog.Any("err", err))
			return
		}
	}

	for _, target := range modInfo.Targets {
		if _, err := db.From(ctx).VersionTarget.Delete().Where(
			versiontarget.VersionID(dbVersion.ID),
			versiontarget.TargetName(target),
		).Exec(schema.SkipSoftDelete(ctx)); err != nil {
			slox.Error(ctx, "failed deleting version target", slog.Any("err", err))
			return
		}
	}

	// For UEPlugin mods, a Windows target is created.
	// However, that happens after the last possible call to this function, therefore we can ignore it

	if err := db.From(ctx).Version.DeleteOneID(dbVersion.ID).Exec(schema.SkipSoftDelete(ctx)); err != nil {
		slox.Error(ctx, "failed deleting version", slog.Any("err", err))
		return
	}

	storage.DeleteMod(ctx, mod.ID, mod.Name, dbVersion.ID)
	for _, target := range modInfo.Targets {
		storage.DeleteModTarget(ctx, mod.ID, mod.Name, dbVersion.ID, target)
	}
}
