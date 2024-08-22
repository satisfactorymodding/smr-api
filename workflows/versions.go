package workflows

import (
	"context"
	"encoding/json"
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
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	mod2 "github.com/satisfactorymodding/smr-api/generated/ent/mod"
	version2 "github.com/satisfactorymodding/smr-api/generated/ent/version"
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

	fatalError := func(ctx workflow.Context, err error, modInfo *validation.ModInfo) error {
		if err != nil {
			err2 := workflow.ExecuteActivity(ctx, removeModActivity, modID, modInfo, uploadID).Get(ctx, nil)
			if err2 != nil {
				slog.Error("failed to remove mod", slog.Any("err", err2))
			}
			return workflow.ExecuteActivity(ctx, storeRedisStateActivity, uploadID, nil, err).Get(ctx, nil)
		}
		return err
	}

	err := workflow.ExecuteActivity(ctx, completeUploadMultipartModActivity, modID, uploadID).Get(ctx, nil)
	if err != nil {
		return fatalError(ctx, err, nil)
	}

	var modInfo *validation.ModInfo
	err = workflow.ExecuteActivity(ctx, extractModInfoActivity, modID, uploadID).Get(ctx, &modInfo)
	if err != nil {
		return fatalError(ctx, err, nil)
	}

	var metadata *string
	err = workflow.ExecuteActivity(ctx, extractMetadataActivity, modID, uploadID, modInfo).Get(ctx, &metadata)
	if err != nil {
		// Do not retry extracting metadata again
		slog.Error("failed to extract metadata", slog.Any("err", err), slog.String("mod_id", modID), slog.String("upload_id", uploadID))
	}

	var fileKey *string
	err = workflow.ExecuteActivity(ctx, renameVersionActivity, modID, uploadID, modInfo.Version).Get(ctx, &fileKey)
	if err != nil {
		return fatalError(ctx, err, modInfo)
	}

	var targetsData *[]modTargetData
	err = workflow.ExecuteActivity(ctx, separateModTargetsActivity, modID, modInfo, fileKey).Get(ctx, &targetsData)
	if err != nil {
		return fatalError(ctx, err, modInfo)
	}

	var dbVersion *ent.Version
	err = workflow.ExecuteActivity(ctx, createVersionInDatabaseActivity, modID, modInfo, fileKey, targetsData, version, metadata).Get(ctx, &dbVersion)
	if err != nil {
		return fatalError(ctx, err, modInfo)
	}

	data := &generated.CreateVersionResponse{
		AutoApproved: shouldAutoApprove(modInfo),
		Version:      (*conv.VersionImpl)(nil).Convert(dbVersion),
	}

	if data.AutoApproved {
		err = workflow.ExecuteActivity(ctx, approveAndPublishModActivity, modID, data.Version.ID).Get(ctx, nil)
		if err != nil {
			return fatalError(ctx, err, modInfo)
		}
	}

	err = workflow.ExecuteActivity(ctx, storeRedisStateActivity, uploadID, data).Get(ctx, nil)
	if err != nil {
		return err
	}

	if !data.AutoApproved {
		var scanSuccess bool
		err = workflow.ExecuteActivity(ctx, scanModOnVirusTotalActivity, modID, data.Version.ID).Get(ctx, &scanSuccess)
		if err != nil {
			return err
		}

		if !scanSuccess {
			return nil
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
		return nil, temporal.NewNonRetryableApplicationError("data.json mod_reference does not match mod reference", "fatal", nil)
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

func createVersionInDatabaseActivity(ctx context.Context, modID string, modInfo *validation.ModInfo, fileKey string, targets []modTargetData, version generated.NewVersion, metadata *string) (*ent.Version, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version", modInfo.Version))

	var dbVersion *ent.Version
	if err := db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
		dbVersion, err = tx.Version.Create().
			SetVersion(modInfo.Version).
			SetGameVersion(modInfo.GameVersion).
			SetChangelog(version.Changelog).
			SetModID(modID).
			SetStability(util.Stability(version.Stability)).
			SetModReference(modInfo.ModReference).
			SetKey(fileKey).
			SetSize(modInfo.Size).
			SetHash(modInfo.Hash).
			SetVersionMajor(int(modInfo.Semver.Major())).
			SetVersionMinor(int(modInfo.Semver.Minor())).
			SetVersionPatch(int(modInfo.Semver.Patch())).
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

		for _, target := range targets {
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

type modTargetData struct {
	TargetName string
	Key        string
	Hash       string
	Size       int64
}

func separateModTargetsActivity(ctx context.Context, modID string, modInfo *validation.ModInfo, fileKey string) ([]modTargetData, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version", modInfo.Version))

	fileData, err := downloadMod(ctx, mod, modInfo.Version)
	if err != nil {
		return nil, err
	}

	if modInfo.Type == validation.MultiTargetUEPlugin {
		targets := make([]modTargetData, 0)

		for _, target := range modInfo.Targets {
			slox.Info(ctx, "separating mod", slog.String("target", target), slog.String("mod", mod.Name), slog.String("version", modInfo.Version))
			key, hash, size, err := storage.SeparateModTarget(ctx, fileData, mod.ID, mod.Name, modInfo.Version, target)
			if err != nil {
				return nil, temporal.NewNonRetryableApplicationError("failed to separate mod", "fatal", err)
			}
			targets = append(targets, modTargetData{
				TargetName: target,
				Key:        key,
				Hash:       hash,
				Size:       size,
			})
		}

		return targets, nil
	}

	// A single Windows target for legacy mod formats
	return []modTargetData{{
		TargetName: "Windows",
		Key:        fileKey,
		Hash:       modInfo.Hash,
		Size:       modInfo.Size,
	}}, nil
}

func renameVersionActivity(ctx context.Context, modID string, uploadID string, version string) (string, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return "", err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("upload_id", uploadID), slog.String("version", version))

	key, err := storage.RenameVersion(ctx, mod.ID, mod.Name, uploadID, version)
	if err != nil {
		return "", temporal.NewNonRetryableApplicationError("failed to upload mod", "fatal", err)
	}
	return key, nil
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

func removeModActivity(ctx context.Context, modID string, modInfo *validation.ModInfo, uploadID string) error {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return err
	}

	// TODO: cleanup file parts if failure happened before completing multipart upload

	_ = storage.DeleteMod(ctx, mod.ID, mod.Name, uploadID)
	if modInfo != nil {
		_ = storage.DeleteMod(ctx, mod.ID, mod.Name, modInfo.Version)
		for _, target := range modInfo.Targets {
			_ = storage.DeleteModTarget(ctx, mod.ID, mod.Name, modInfo.Version, target)
		}
	}
	return nil
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

func shouldAutoApprove(modInfo *validation.ModInfo) bool {
	if viper.GetBool("skip-virus-check") {
		return true
	}

	for _, obj := range modInfo.Objects {
		if obj.Type != "pak" {
			return false
		}
	}

	return true
}
