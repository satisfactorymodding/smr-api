package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Vilsol/slox"
	"github.com/pkg/errors"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/validation"
)

func UpdateModDataFromStorage(ctx context.Context, modID string, versionID string, metadata bool) error {
	// perform task
	slox.Info(ctx, "Updating DB for mod version with metadata", slog.String("mod", modID), slog.String("version", versionID), slog.Bool("metadata", metadata))

	version := postgres.GetVersion(ctx, versionID)
	link := storage.GenerateDownloadLink(version.Key)

	response, _ := http.Get(link)

	fileData, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	mod := postgres.GetModByID(ctx, modID)

	if mod == nil {
		return errors.New("mod not found")
	}

	info, err := validation.ExtractModInfo(ctx, fileData, metadata, false, mod.ModReference)
	if err != nil {
		slox.Warn(ctx, "failed updating mod, likely outdated", slog.Any("err", err), slog.String("version", versionID))
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
			slox.Error(ctx, "failed serializing", slog.Any("err", err), slog.String("version", versionID))
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
