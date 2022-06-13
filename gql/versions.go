package gql

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/integrations"
	"github.com/satisfactorymodding/smr-api/redis/jobs"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
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
	fileData, err := ioutil.ReadAll(modFile)

	if err != nil {
		storage.DeleteMod(ctx, mod.ID, mod.Name, versionID)
		return nil, errors.Wrap(err, "failed reading mod file")
	}

	modInfo, err := validation.ExtractModInfo(ctx, fileData, true, true, mod.ModReference)

	if err != nil {
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
		//Size:         &modInfo.Size,
		//Hash:         &modInfo.Hash,
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
		log.Ctx(ctx).Err(err).Msgf("[%s] failed serializing", dbVersion.ID)
	} else {
		metadata := string(jsonData)
		dbVersion.Metadata = &metadata
		postgres.Save(ctx, &dbVersion)
	}

	// TODO Validate mod files
	key, err := storage.SeparateMod(ctx, fileData, mod.ID, mod.Name, versionID, modInfo.Version)

	if err != nil {
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
		storage.DeleteCombinedMod(ctx, mod.ID, mod.Name, dbVersion.ID)
		return nil, errors.New("failed to upload mod")
	}

	dbVersion.Key = key
	postgres.Save(ctx, &dbVersion)
	postgres.Save(ctx, &mod)

	if autoApproved {
		mod := postgres.GetModByID(ctx, dbVersion.ModID)
		now := time.Now()
		mod.LastVersionDate = &now
		postgres.Save(ctx, &mod)

		go integrations.NewVersion(util.ReWrapCtx(ctx), dbVersion)
		storage.DeleteCombinedMod(ctx, mod.ID, mod.Name, dbVersion.ID)
	} else {
		l.Info().Msg("Submitting version job for virus scan")
		jobs.SubmitJobScanModOnVirusTotalTask(ctx, mod.ID, dbVersion.ID, true)
	}

	return &generated.CreateVersionResponse{
		AutoApproved: autoApproved,
		Version:      DBVersionToGenerated(dbVersion),
	}, nil
}
