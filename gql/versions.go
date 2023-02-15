package gql

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/redis/jobs"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"
)

func FinalizeVersionUploadAsync(ctx context.Context, mod *postgres.Mod, versionID string, version generated.NewVersion) (*generated.CreateVersionResponse, error) {
	l := log.With().Str("mod_id", mod.ID).Str("version_id", versionID).Logger()

	l.Info().Msg("Creating multipart upload")
	success, _ := storage.CompleteUploadMultipartMod(ctx, mod.ID, mod.Name, versionID)

	if !success {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, errors.New("failed uploading mod")
	}

	modFile, err := storage.GetMod(mod.ID, mod.Name, versionID)
	if err != nil {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, err
	}

	// TODO Optimize
	fileData, err := io.ReadAll(modFile)
	if err != nil {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, errors.Wrap(err, "failed reading mod file")
	}

	modInfo, err := validation.ExtractModInfo(ctx, fileData, true, true, mod.ModReference)
	if err != nil {
		spew.Dump(err)
		l.Err(err).Msg("failed extracting mod info")
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, err
	}

	if modInfo.ModReference != mod.ModReference {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, errors.New("data.json mod_reference does not match mod reference")
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
		log.Err(err).Msgf("[%s] failed serializing", dbVersion.ID)
	} else {
		metadata := string(jsonData)
		dbVersion.Metadata = &metadata
		postgres.Save(ctx, &dbVersion)
	}

	separated := storage.SeparateMod(ctx, fileData, mod.ID, mod.Name, dbVersion.ID, modInfo.Version)

	if !separated {
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

		postgres.DeleteForced(ctx, &dbVersion)
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)

		for _, dbModArch := range dbVersion.Arch {
			postgres.DeleteForced(ctx, &dbModArch)
		}
		return nil, errors.New("failed to upload mod")
	}

	dbModArch := postgres.GetModArchByPlatform(ctx, dbVersion.ID, "WindowsNoEditor")

	if dbModArch == nil {
		dbVersion.Key = ""
		dbVersion.Hash = nil
		dbVersion.Size = nil
	} else {
		dbVersion.Key = dbModArch.Key
		dbVersion.Hash = &dbModArch.Hash
		dbVersion.Size = &dbModArch.Size
	}

	postgres.Save(ctx, &dbVersion)
	postgres.Save(ctx, &mod)

	storage.DeleteVersion(ctx, mod.ID, mod.Name, versionID)

	if autoApproved {
		mod := postgres.GetModByID(ctx, dbVersion.ModID)
		now := time.Now()
		mod.LastVersionDate = &now
		postgres.Save(ctx, &mod)

		go integrations.NewVersion(util.ReWrapCtx(ctx), dbVersion)
	} else {
		l.Info().Msg("Submitting version job for virus scan")
		jobs.SubmitJobScanModOnVirusTotalTask(ctx, mod.ID, dbVersion.ID, true)
	}

	return &generated.CreateVersionResponse{
		AutoApproved: autoApproved,
		Version:      DBVersionToGenerated(dbVersion),
	}, nil
}
