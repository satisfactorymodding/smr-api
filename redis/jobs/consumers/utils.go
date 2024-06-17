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

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/validation"
)

func UpdateModDataFromStorage(ctx context.Context, modID string, versionID string, metadata bool) error {
	// perform task
	slox.Info(ctx, "Updating DB for mod version with metadata", slog.String("mod", modID), slog.String("version", versionID), slog.Bool("metadata", metadata))

	version, err := db.From(ctx).Version.Get(ctx, versionID)
	if err != nil {
		return err
	}

	link := storage.GenerateDownloadLink(version.Key)

	response, _ := http.Get(link)

	fileData, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	if err != nil {
		return err
	}

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
		if err := db.From(ctx).VersionDependency.Create().
			SetVersionID(version.ID).
			SetModID(depModID).
			SetCondition(condition).
			SetOptional(false).
			OnConflict().
			DoNothing().
			Exec(ctx); err != nil {
			return err
		}
	}

	for depModID, condition := range info.OptionalDependencies {
		if err := db.From(ctx).VersionDependency.Create().
			SetVersionID(version.ID).
			SetModID(depModID).
			SetCondition(condition).
			SetOptional(true).
			OnConflict().
			DoNothing().
			Exec(ctx); err != nil {
			return err
		}
	}

	if metadata {
		jsonData, err := json.Marshal(info.Metadata)
		if err != nil {
			slox.Error(ctx, "failed serializing", slog.Any("err", err), slog.String("version", versionID))
		} else {
			if err := version.Update().SetMetadata(string(jsonData)).Exec(ctx); err != nil {
				return err
			}
		}
	}

	versionMajor := int(info.Semver.Major())
	versionMinor := int(info.Semver.Minor())
	versionPatch := int(info.Semver.Patch())

	return version.Update().
		SetSize(info.Size).
		SetHash(info.Hash).
		SetVersionMajor(versionMajor).
		SetVersionMinor(versionMinor).
		SetVersionPatch(versionPatch).
		SetModReference(info.ModReference).
		SetSmlVersion(info.SMLVersion).
		Exec(ctx)
}
