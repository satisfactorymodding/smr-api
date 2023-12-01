package gql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log/slog"
	"time"

	"github.com/Vilsol/slox"
	"github.com/pkg/errors"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/redis/jobs"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"
)

func FinalizeVersionUploadAsync(ctx context.Context, mod *postgres.Mod, versionID string, version generated.NewVersion) (*generated.CreateVersionResponse, error) {
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

	dbVersion := &postgres.Version{
		Version:      modInfo.Version,
		SMLVersion:   modInfo.SMLVersion,
		Changelog:    version.Changelog,
		ModID:        mod.ID,
		Stability:    string(version.Stability),
		ModReference: &modInfo.ModReference,
		Size:         &modInfo.Size,
		Hash:         &modInfo.Hash,
		VersionMajor: &versionMajor,
		VersionMinor: &versionMinor,
		VersionPatch: &versionPatch,
	}

	autoApproved := true
	for _, obj := range modInfo.Objects {
		if obj.Type != "pak" {
			autoApproved = false
			break
		}
	}

	autoApproved = autoApproved || viper.GetBool("skip-virus-check")

	dbVersion.Approved = autoApproved

	err = postgres.CreateVersion(ctx, dbVersion)

	if err != nil {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, err
	}

	for modID, condition := range modInfo.Dependencies {
		dependency := postgres.VersionDependency{
			VersionID: dbVersion.ID,
			ModID:     modID,
			Condition: condition,
			Optional:  false,
		}

		postgres.Save(ctx, &dependency)
	}

	for modID, condition := range modInfo.OptionalDependencies {
		dependency := postgres.VersionDependency{
			VersionID: dbVersion.ID,
			ModID:     modID,
			Condition: condition,
			Optional:  true,
		}

		postgres.Save(ctx, &dependency)
	}

	jsonData, err := json.Marshal(modInfo.Metadata)
	if err != nil {
		slox.Error(ctx, "failed serializing", slog.Any("err", err), slog.String("version_id", dbVersion.ID))
	} else {
		metadata := string(jsonData)
		dbVersion.Metadata = &metadata
		postgres.Save(ctx, &dbVersion)
	}

	if modInfo.Type == validation.MultiTargetUEPlugin {
		targets := make([]*postgres.VersionTarget, 0)

		for _, target := range modInfo.Targets {
			dbVersionTarget := &postgres.VersionTarget{
				VersionID:  dbVersion.ID,
				TargetName: target,
			}

			postgres.Save(ctx, dbVersionTarget)

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

			target.Key = key
			target.Hash = hash
			target.Size = size

			postgres.Save(ctx, target)
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
		dbVersionTarget := &postgres.VersionTarget{
			VersionID:  dbVersion.ID,
			TargetName: "Windows",
			Key:        key,
			Hash:       *dbVersion.Hash,
			Size:       *dbVersion.Size,
		}

		postgres.Save(ctx, dbVersionTarget)
	}

	dbVersion.Key = key
	postgres.Save(ctx, &dbVersion)
	postgres.Save(ctx, &mod)

	if autoApproved {
		mod := postgres.GetModByID(ctx, dbVersion.ModID)
		now := time.Now()
		mod.LastVersionDate = &now
		postgres.Save(ctx, &mod)

		go integrations.NewVersion(db.ReWrapCtx(ctx), dbVersion)
	} else {
		slox.Info(ctx, "Submitting version job for virus scan")
		jobs.SubmitJobScanModOnVirusTotalTask(ctx, mod.ID, dbVersion.ID, true)
	}

	return &generated.CreateVersionResponse{
		AutoApproved: autoApproved,
		Version:      DBVersionToGenerated(dbVersion),
	}, nil
}

func removeMod(ctx context.Context, modInfo *validation.ModInfo, mod *postgres.Mod, dbVersion *postgres.Version) {
	for modID, condition := range modInfo.Dependencies {
		dependency := postgres.VersionDependency{
			VersionID: dbVersion.ID,
			ModID:     modID,
			Condition: condition,
			Optional:  false,
		}

		postgres.DeleteForced(ctx, &dependency)
	}

	for modID, condition := range modInfo.OptionalDependencies {
		dependency := postgres.VersionDependency{
			VersionID: dbVersion.ID,
			ModID:     modID,
			Condition: condition,
			Optional:  true,
		}

		postgres.DeleteForced(ctx, &dependency)
	}

	for _, target := range modInfo.Targets {
		dbVersionTarget := postgres.VersionTarget{
			VersionID:  dbVersion.ID,
			TargetName: target,
		}

		postgres.DeleteForced(ctx, &dbVersionTarget)
	}

	// For UEPlugin mods, a Windows target is created.
	// However, that happens after the last possible call to this function, therefore we can ignore it

	postgres.DeleteForced(ctx, &dbVersion)

	storage.DeleteMod(ctx, mod.ID, mod.Name, dbVersion.ID)
	for _, target := range modInfo.Targets {
		storage.DeleteModTarget(ctx, mod.ID, mod.Name, dbVersion.ID, target)
	}
}
