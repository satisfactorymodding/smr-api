package consumers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/validation"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
)

func UpdateModDataFromStorage(ctx context.Context, modID string, versionID string, metadata bool) error {
	// perform task
	log.Info().Msgf("[%s] Updating DB for mod %s version with metadata: %v", versionID, modID, metadata)

	version := postgres.GetVersion(ctx, versionID)
	link := storage.GenerateDownloadLink(version.Key)

	response, _ := http.Get(link)

	fileData, err := io.ReadAll(response.Body)

	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	mod := postgres.GetModByID(ctx, modID)

	if mod == nil {
		return errors.New("mod not found")
	}

	info, err := validation.ExtractModInfo(ctx, fileData, metadata, false, mod.ModReference)

	if err != nil {
		log.Warn().Err(err).Msgf("[%s] Failed updating mod, likely outdated", versionID)
		// Outdated version
		return nil
	}

	for depModID, condition := range info.Dependencies {
		dependency := postgres.VersionDependency{
			VersionID: version.ID,
			ModID:     depModID,
			Condition: condition,
			Optional:  false,
		}

		postgres.Save(ctx, &dependency)
	}

	for depModID, condition := range info.OptionalDependencies {
		dependency := postgres.VersionDependency{
			VersionID: version.ID,
			ModID:     depModID,
			Condition: condition,
			Optional:  true,
		}

		postgres.Save(ctx, &dependency)
	}

	if metadata {
		jsonData, err := json.Marshal(info.Metadata)
		if err != nil {
			log.Err(err).Msgf("[%s] failed serializing", versionID)
		} else {
			metadata := string(jsonData)
			version.Metadata = &metadata
		}
	}

	versionMajor := int(info.Semver.Major())
	versionMinor := int(info.Semver.Minor())
	versionPatch := int(info.Semver.Patch())

	version.Size = &info.Size
	version.Hash = &info.Hash
	version.VersionMajor = &versionMajor
	version.VersionMinor = &versionMinor
	version.VersionPatch = &versionPatch

	version.ModReference = &info.ModReference
	version.SMLVersion = info.SMLVersion
	postgres.Save(ctx, &version)

	return nil
}
