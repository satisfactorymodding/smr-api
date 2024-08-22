package workflows

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Vilsol/slox"
	"github.com/pkg/errors"
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

func FinalizeVersionUploadWorkflow(ctx workflow.Context, modID string, versionID string, version generated.NewVersion) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{"UnrecoverableError"},
		},
	})

	err := workflow.ExecuteActivity(ctx, completeUploadMultipartModActivity, modID, versionID).Get(ctx, nil)
	if err != nil {
		return err
	}

	var modInfo *validation.ModInfo
	err = workflow.ExecuteActivity(ctx, extractModInfo, modID, versionID).Get(ctx, &modInfo)
	if err != nil {
		return err
	}

	var dbVersion *ent.Version
	err = workflow.ExecuteActivity(ctx, createVersionInDatabaseActivity, modID, versionID, modInfo, version).Get(ctx, &dbVersion)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, separateModTargetsActivity, modID, versionID, modInfo, dbVersion).Get(ctx, nil)
	if err != nil {
		return err
	}

	var data *generated.CreateVersionResponse
	err = workflow.ExecuteActivity(ctx, finalizeVersionUploadActivity, modID, versionID, modInfo, dbVersion).Get(ctx, &data)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, storeRedisStateActivity, versionID, data).Get(ctx, &data)
	if err != nil {
		return err
	}

	if !data.AutoApproved {
		err = workflow.ExecuteActivity(ctx, scanModOnVirusTotalActivity, modID, data.Version.ID, true).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func storeRedisStateActivity(ctx context.Context, versionID string, data *generated.CreateVersionResponse, err error) error {
	ctx = slox.With(ctx, slog.String("version_id", versionID))

	if err2 := redis.StoreVersionUploadState(versionID, data, err); err2 != nil {
		slox.Error(ctx, "error storing redis state", slog.Any("err", err2))
		return err2
	}

	return nil
}

func completeUploadMultipartModActivity(ctx context.Context, modID string, versionID string) error {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version_id", versionID))

	slox.Info(ctx, "Completing multipart upload")
	success, _ := storage.CompleteUploadMultipartMod(ctx, mod.ID, mod.Name, versionID)

	if !success {
		return errors.New("failed uploading mod")
	}

	return nil
}

func extractModInfo(ctx context.Context, modID string, versionID string) (*validation.ModInfo, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version_id", versionID))

	fileData, err := downloadMod(mod, versionID)
	if err != nil {
		return nil, err
	}

	modInfo, err := validation.ExtractModInfo(ctx, fileData, true, true, mod.ModReference)
	if err != nil {
		return nil, fmt.Errorf("failed extracting mod info: %w", err)
	}

	if modInfo.ModReference != mod.ModReference {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, UnrecoverableError{errors.New("data.json mod_reference does not match mod reference")}
	}

	if modInfo.Type == validation.DataJSON {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, UnrecoverableError{errors.New("data.json mods are obsolete and not allowed")}
	}

	if modInfo.Type == validation.MultiTargetUEPlugin && !util.FlagEnabled(util.FeatureFlagAllowMultiTargetUpload) {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, UnrecoverableError{errors.New("multi-target mods are not allowed")}
	}

	count, err := db.From(ctx).Version.Query().
		Where(version2.ModID(mod.ID), version2.Version(modInfo.Version)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, UnrecoverableError{errors.New("this mod already has a version with this name")}
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
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		timeToWait := time.Until(versions[0].CreatedAt.Add(time.Hour * 24)).Minutes()
		return nil, UnrecoverableError{fmt.Errorf("please wait %.0f minutes to post another version", timeToWait)}
	}

	return modInfo, nil
}

func createVersionInDatabaseActivity(ctx context.Context, modID string, versionID string, modInfo *validation.ModInfo, version generated.NewVersion) (*ent.Version, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version_id", versionID))

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
			SetNillableMetadata(modInfo.MetadataJSON).
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
		storage.DeleteMod(ctx, modID, mod.Name, versionID)
		return nil, UnrecoverableError{err}
	}

	return dbVersion, nil
}

func separateModTargetsActivity(ctx context.Context, modID string, versionID string, modInfo *validation.ModInfo, dbVersion *ent.Version) error {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version_id", versionID))

	fileData, err := downloadMod(mod, versionID)
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
				storage.DeleteMod(ctx, modID, mod.Name, versionID)
				return UnrecoverableError{err}
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
				removeMod(ctx, modInfo, mod, dbVersion)
				return UnrecoverableError{err}
			}
		}

		if !separateSuccess {
			removeMod(ctx, modInfo, mod, dbVersion)
			return UnrecoverableError{errors.New("failed to separate mod")}
		}
	}

	return nil
}

func finalizeVersionUploadActivity(ctx context.Context, modID string, versionID string, modInfo *validation.ModInfo, dbVersion *ent.Version) (*generated.CreateVersionResponse, error) {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return nil, err
	}

	ctx = slox.With(ctx, slog.String("mod_id", mod.ID), slog.String("version_id", versionID))

	success, key := storage.RenameVersion(ctx, mod.ID, mod.Name, versionID, modInfo.Version)

	if !success {
		removeMod(ctx, modInfo, mod, dbVersion)
		return nil, UnrecoverableError{errors.New("failed to upload mod")}
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
			return nil, UnrecoverableError{err}
		}
	}

	if _, err := db.From(ctx).Version.Update().SetKey(key).Where(version2.ID(dbVersion.ID)).Save(ctx); err != nil {
		removeMod(ctx, modInfo, mod, dbVersion)
		return nil, UnrecoverableError{err}
	}

	if dbVersion.Approved {
		if _, err := db.From(ctx).Mod.Update().SetLastVersionDate(time.Now()).Where(mod2.ID(mod.ID)).Save(ctx); err != nil {
			removeMod(ctx, modInfo, mod, dbVersion)
			return nil, UnrecoverableError{err}
		}

		go integrations.NewVersion(db.ReWrapCtx(ctx), dbVersion)
	} else {
		slox.Info(ctx, "version will require virus scan")
	}

	return &generated.CreateVersionResponse{
		AutoApproved: dbVersion.Approved,
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

func downloadMod(mod *ent.Mod, versionID string) ([]byte, error) {
	modFile, err := storage.GetMod(mod.ID, mod.Name, versionID)
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
