package workflows

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Vilsol/slox"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/db/schema"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	mod2 "github.com/satisfactorymodding/smr-api/generated/ent/mod"
	version2 "github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiondependency"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiontarget"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"
)

func FinalizeVersionUploadWorkflow(ctx workflow.Context, modID string, uploadID string, version generated.NewVersion) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 10,
		},
	})

	storeIfFatal := func(ctx workflow.Context, err error, data *generated.CreateVersionResponse) error {
		if err != nil {
			var appError *temporal.ApplicationError
			if errors.As(err, &appError) && appError.NonRetryable() {
				return workflow.ExecuteActivity(ctx, storeRedisStateActivity, uploadID, data).Get(ctx, nil)
			}
		}
		return err
	}

	err := workflow.ExecuteActivity(ctx, completeUploadMultipartModActivity, modID, uploadID).Get(ctx, nil)
	if err != nil {
		return storeIfFatal(ctx, err, nil)
	}

	var modInfo *validation.ModInfo
	err = workflow.ExecuteActivity(ctx, extractModInfoActivity, modID, uploadID).Get(ctx, &modInfo)
	if err != nil {
		return storeIfFatal(ctx, err, nil)
	}

	var metadata *string
	err = workflow.ExecuteActivity(ctx, extractMetadataActivity, modID, uploadID, modInfo).Get(ctx, &metadata)
	if err != nil {
		// Do not retry extracting metadata again
		slog.Error("failed to extract metadata", slog.Any("err", err), slog.String("mod_id", modID), slog.String("upload_id", uploadID))
	}

	var dbVersion *ent.Version
	err = workflow.ExecuteActivity(ctx, createVersionInDatabaseActivity, modID, uploadID, modInfo, version, metadata).Get(ctx, &dbVersion)
	if err != nil {
		return storeIfFatal(ctx, err, nil)
	}

	err = workflow.ExecuteActivity(ctx, separateModTargetsActivity, modID, uploadID, modInfo, dbVersion).Get(ctx, nil)
	if err != nil {
		return storeIfFatal(ctx, err, nil)
	}

	var data *generated.CreateVersionResponse
	err = workflow.ExecuteActivity(ctx, finalizeVersionUploadActivity, modID, uploadID, modInfo, dbVersion).Get(ctx, &data)
	if err != nil {
		return storeIfFatal(ctx, err, nil)
	}

	if data.AutoApproved {
		err = workflow.ExecuteActivity(ctx, approveAndPublishModActivity, modID, data.Version.ID).Get(ctx, nil)
		if err != nil {
			return storeIfFatal(ctx, err, nil)
		}
	}

	err = workflow.ExecuteActivity(ctx, storeRedisStateActivity, uploadID, data).Get(ctx, nil)
	if err != nil {
		return err
	}

	if !data.AutoApproved {
		err = workflow.ExecuteActivity(ctx, scanModOnVirusTotalActivity, modID, data.Version.ID).Get(ctx, nil)
		if err != nil {
			return err
		}

		err = workflow.ExecuteActivity(ctx, approveAndPublishModActivity, modID, data.Version.ID).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func completeUploadMultipartModActivity(ctx context.Context, modID string, uploadID string) error {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", uploadID))

	slox.Info(ctx, "Completing multipart upload")
	_, err = storage.CompleteUploadMultipartMod(ctx, mod.ID, mod.Name, uploadID)
	return err
}

func extractModInfoActivity(ctx context.Context, modID string, uploadID string) (*validation.ModInfo, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", uploadID))

	fileData, err := downloadMod(ctx, mod, uploadID)
	if err != nil {
		return nil, err
	}

	modInfo, err := validation.ExtractModInfo(ctx, fileData, true, mod.ModReference)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("failed extracting mod info", "fatal", err)
	}

	if modInfo.ModReference != mod.ModReference {
		_ = storage.DeleteMod(ctx, mod.ID, mod.Name, uploadID)
		return nil, temporal.NewNonRetryableApplicationError("data.json mod_reference does not match mod reference", "fatal", nil)
	}

	if modInfo.Type == validation.DataJSON {
		_ = storage.DeleteMod(ctx, mod.ID, mod.Name, uploadID)
		return nil, temporal.NewNonRetryableApplicationError("data.json mods are obsolete and not allowed", "fatal", nil)
	}

	if modInfo.Type == validation.MultiTargetUEPlugin && !util.FlagEnabled(util.FeatureFlagAllowMultiTargetUpload) {
		_ = storage.DeleteMod(ctx, mod.ID, mod.Name, uploadID)
		return nil, temporal.NewNonRetryableApplicationError("multi-target mods are not allowed", "fatal", nil)
	}

	count, err := db.From(ctx).Version.Query().
		Where(version2.ModID(mod.ID), version2.Version(modInfo.Version)).
		Count(ctx)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("database error", "fatal", err)
	}

	if count > 0 {
		_ = storage.DeleteMod(ctx, mod.ID, mod.Name, uploadID)
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
		_ = storage.DeleteMod(ctx, mod.ID, mod.Name, uploadID)
		timeToWait := time.Until(versions[0].CreatedAt.Add(time.Hour * 24)).Minutes()
		return nil, temporal.NewNonRetryableApplicationError(fmt.Sprintf("please wait %.0f minutes to post another version", timeToWait), "fatal", nil)
	}

	return modInfo, nil
}

func extractMetadataActivity(ctx context.Context, modID string, uploadID string, modInfo *validation.ModInfo) (*string, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", uploadID))

	fileData, err := downloadMod(ctx, mod, uploadID)
	if err != nil {
		return nil, err
	}

	metadata, err := validation.ExtractMetadata(ctx, fileData, modInfo.GameVersion, modInfo.ModReference)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed json serialization: %w", err)
	}

	return util.Ptr(string(jsonData)), nil
}

func createVersionInDatabaseActivity(ctx context.Context, modID string, uploadID string, modInfo *validation.ModInfo, version generated.NewVersion, metadata *string) (*ent.Version, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", uploadID))

	autoApproved := true
	for _, obj := range modInfo.Objects {
		if obj.Type != "pak" {
			autoApproved = false
			break
		}
	}

	autoApproved = autoApproved || viper.GetBool("skip-virus-check")

	var dbVersion *ent.Version
	if err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		dbVersion, err = tx.Version.Create().
			SetVersion(modInfo.Version).
			SetGameVersion(modInfo.GameVersion).
			SetChangelog(version.Changelog).
			SetModID(modID).
			SetStability(util.Stability(version.Stability)).
			SetModReference(modInfo.ModReference).
			SetSize(modInfo.Size).
			SetHash(modInfo.Hash).
			SetVersionMajor(int(modInfo.Semver.Major())).
			SetVersionMinor(int(modInfo.Semver.Minor())).
			SetVersionPatch(int(modInfo.Semver.Patch())).
			SetApproved(autoApproved).
			SetNillableMetadata(metadata).
			Save(ctx)
		if err != nil {
			return err
		}

		for modReference, condition := range modInfo.Dependencies {
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

		for modReference, condition := range modInfo.OptionalDependencies {
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

		return nil
	}, nil); err != nil {
		_ = storage.DeleteMod(ctx, modID, mod.Name, uploadID)
		return nil, temporal.NewNonRetryableApplicationError("database error", "fatal", err)
	}

	return dbVersion, nil
}

func separateModTargetsActivity(ctx context.Context, modID string, uploadID string, modInfo *validation.ModInfo, dbVersion *ent.Version) error {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", uploadID))

	fileData, err := downloadMod(ctx, mod, uploadID)
	if err != nil {
		return err
	}

	if modInfo.Type == validation.MultiTargetUEPlugin {
		targets := make([]*ent.VersionTarget, 0)

		for _, target := range modInfo.Targets {
			dbVersionTarget, err := db.From(ctx).VersionTarget.Create().
				SetVersionID(dbVersion.ID).
				SetTargetName(target).
				Save(ctx)
			if err != nil {
				_ = storage.DeleteMod(ctx, modID, mod.Name, uploadID)
				return temporal.NewNonRetryableApplicationError("database error", "fatal", err)
			}

			targets = append(targets, dbVersionTarget)
		}

		var separateError error
		for _, target := range targets {
			slox.Info(ctx, "separating mod", slog.String("target", target.TargetName), slog.String("mod", mod.Name), slog.String("version", dbVersion.Version))
			key, hash, size, err := storage.SeparateModTarget(ctx, fileData, mod.ID, mod.Name, dbVersion.Version, target.TargetName)
			if err != nil {
				separateError = err
				break
			}

			if _, err := target.Update().SetKey(key).SetHash(hash).SetSize(size).Save(ctx); err != nil {
				removeMod(ctx, modInfo, mod, dbVersion)
				return temporal.NewNonRetryableApplicationError("database error", "fatal", err)
			}
		}

		if separateError != nil {
			removeMod(ctx, modInfo, mod, dbVersion)
			return temporal.NewNonRetryableApplicationError("failed to separate mod", "fatal", err)
		}
	}

	return nil
}

func finalizeVersionUploadActivity(ctx context.Context, modID string, uploadID string, modInfo *validation.ModInfo, dbVersion *ent.Version) (*generated.CreateVersionResponse, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", uploadID))

	key, err := storage.RenameVersion(ctx, mod.ID, mod.Name, uploadID, modInfo.Version)
	if err != nil {
		removeMod(ctx, modInfo, mod, dbVersion)
		return nil, temporal.NewNonRetryableApplicationError("failed to upload mod", "fatal", err)
	}

	if modInfo.Type == validation.UEPlugin {
		if _, err := db.From(ctx).VersionTarget.Create().
			SetVersionID(dbVersion.ID).
			SetTargetName("Windows").
			SetKey(key).
			SetHash(dbVersion.Hash).
			SetSize(dbVersion.Size).
			Save(ctx); err != nil {
			removeMod(ctx, modInfo, mod, dbVersion)
			return nil, temporal.NewNonRetryableApplicationError("database error", "fatal", err)
		}
	}

	if _, err := db.From(ctx).Version.Update().SetKey(key).Where(version2.ID(dbVersion.ID)).Save(ctx); err != nil {
		removeMod(ctx, modInfo, mod, dbVersion)
		return nil, temporal.NewNonRetryableApplicationError("database error", "fatal", err)
	}

	return &generated.CreateVersionResponse{
		AutoApproved: dbVersion.Approved,
		Version:      (*conv.VersionImpl)(nil).Convert(dbVersion),
	}, nil
}

func approveAndPublishModActivity(ctx context.Context, modID string, versionID string) error {
	slox.Info(ctx, "approving mod", slog.String("mod", modID), slog.String("version", versionID))

	version, err := db.From(ctx).Version.Get(ctx, versionID)
	if err != nil {
		return err
	}

	if err := version.Update().SetApproved(true).Exec(ctx); err != nil {
		return err
	}

	if err := db.From(ctx).Mod.UpdateOneID(modID).SetLastVersionDate(time.Now()).Exec(ctx); err != nil {
		return err
	}

	go integrations.NewVersion(db.ReWrapCtx(ctx), version)

	return nil
}

func storeRedisStateActivity(ctx context.Context, uploadID string, data *generated.CreateVersionResponse, err error) error {
	ctx = slox.With(ctx, slog.String("upload_id", uploadID))

	if err2 := redis.StoreVersionUploadState(uploadID, data, err); err2 != nil {
		slox.Error(ctx, "error storing redis state", slog.Any("err", err2))
		return err2
	}

	return nil
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

	_ = storage.DeleteMod(ctx, mod.ID, mod.Name, dbVersion.ID)
	for _, target := range modInfo.Targets {
		_ = storage.DeleteModTarget(ctx, mod.ID, mod.Name, dbVersion.ID, target)
	}
}

func downloadMod(ctx context.Context, mod *ent.Mod, uploadID string) ([]byte, error) {
	modFile, err := storage.GetMod(ctx, mod.ID, mod.Name, uploadID)
	if err != nil {
		return nil, fmt.Errorf("failed getting mod: %w", err)
	}

	// TODO Optimize
	fileData, err := io.ReadAll(modFile)
	if err != nil {
		return nil, fmt.Errorf("failed reading mod file: %w", err)
	}

	return fileData, nil
}
