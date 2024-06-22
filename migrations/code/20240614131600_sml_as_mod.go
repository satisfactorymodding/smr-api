package code

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/Masterminds/semver/v3"
	"github.com/lab259/go-migration"
	"github.com/pkg/errors"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/migrations/utils"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/validation"
)

func init() {
	migration.NewCodeMigration(
		func(_ interface{}) error {
			ctx, err := db.WithDB(context.Background())
			if err != nil {
				return err
			}

			err = uploadAllSMLVersions(ctx)
			if err != nil {
				return fmt.Errorf("failed to upload all SML versions: %w", err)
			}

			// Add the game version
			err = utils.ReindexAllModFiles(ctx, false, nil, nil)
			if err != nil {
				return fmt.Errorf("failed to reindex all mod files: %w", err)
			}

			return nil
		},
	)
}

func uploadAllSMLVersions(ctx context.Context) error {
	smlModQuery := db.From(ctx).Mod.Query().Where(mod.ModReference("SML"))
	smlModQuery = smlModQuery.WithVersions(func(q *ent.VersionQuery) {
		q.WithTargets()
	})

	smlMod, err := smlModQuery.First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get SML mod: %w", err)
	}

	for _, version := range smlMod.Edges.Versions {
		factoryGameVersion := version.GameVersion[2:] // Strip the >=

		var archive []byte
		extractInfo := true

		semVersion, err := semver.NewVersion(version.Version)
		if err != nil {
			return fmt.Errorf("failed to parse SML version %s: %w", version.Version, err)
		}

		switch semVersion.Major() {
		case 1:
			archive, err = createSML1Archive(version)
			if err != nil {
				return fmt.Errorf("failed to create SML1 archive: %w", err)
			}
			extractInfo = false // Not a valid SMR archive
		case 2:
			archive, err = createSML2Archive(version)
			if err != nil {
				return fmt.Errorf("failed to create SML2 archive: %w", err)
			}
			extractInfo = false // Not a valid SMR archive
		case 3:
			archive, err = createSML3Archive(version.Edges.Targets, processSMLUplugin(factoryGameVersion))
			if err != nil {
				return fmt.Errorf("failed to create SML3 archive: %w", err)
			}
		default:
			return fmt.Errorf("unknown SML version %s", version.Version)
		}
		err = uploadSMLVersion(ctx, version, archive, extractInfo)
		if err != nil {
			return fmt.Errorf("failed to upload SML %s: %w", version.Version, err)
		}
	}

	return nil
}

func uploadSMLVersion(ctx context.Context, version *ent.Version, archive []byte, extractInfo bool) error {
	ok, _ := storage.StartUploadMultipartMod(ctx, version.ModID, "SML", version.ID)
	if !ok {
		return errors.New("failed to start upload")
	}

	ok, _ = storage.UploadMultipartMod(ctx, version.ModID, "SML", version.ID, 1, bytes.NewReader(archive))
	if !ok {
		return errors.New("failed to upload")
	}

	ok, _ = storage.CompleteUploadMultipartMod(ctx, version.ModID, "SML", version.ID)
	if !ok {
		return errors.New("failed to complete upload")
	}

	for _, target := range version.Edges.Targets {
		success, key, hash, size := storage.SeparateModTarget(ctx, archive, version.ModID, "SML", version.Version, target.TargetName)
		if !success {
			return errors.New("failed to separate target")
		}

		err := db.From(ctx).VersionTarget.UpdateOneID(target.ID).
			SetKey(key).
			SetHash(hash).
			SetSize(size).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to update target %s: %w", target.TargetName, err)
		}
	}

	ok, key := storage.RenameVersion(ctx, version.ModID, "SML", version.ID, version.Version)
	if !ok {
		return errors.New("failed to rename version")
	}

	if !extractInfo {
		return nil
	}

	modInfo, err := validation.ExtractModInfo(ctx, archive, false, false, "SML")
	if err != nil {
		return fmt.Errorf("failed to extract mod info: %w", err)
	}

	err = db.From(ctx).Version.UpdateOneID(version.ID).
		SetKey(key).
		SetHash(modInfo.Hash).
		SetSize(modInfo.Size).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update version: %w", err)
	}

	return nil
}

func createArchiveFromURLs(urls ...string) ([]byte, error) {
	archive := bytes.NewBuffer([]byte{})
	archiveWriter := zip.NewWriter(archive)

	for _, u := range urls {
		err := addURLToArchive(archiveWriter, u)
		if err != nil {
			return nil, fmt.Errorf("failed to add URL to archive: %w", err)
		}
	}

	err := archiveWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close archive: %w", err)
	}

	return archive.Bytes(), nil
}

func addURLToArchive(zip *zip.Writer, u string) error {
	parsed, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	filename := path.Base(parsed.Path)

	resp, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	file, err := zip.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func createSML1Archive(version *ent.Version) ([]byte, error) {
	return createArchiveFromURLs(
		fmt.Sprintf("https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v%s/xinput1_3.dll", version.Version),
	)
}

func createSML2Archive(version *ent.Version) ([]byte, error) {
	return createArchiveFromURLs(
		fmt.Sprintf("https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v%s/SML.pak", version.Version),
		fmt.Sprintf("https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v%s/UE4-SML-Win64-Shipping.dll", version.Version),
	)
}

func createSML3Archive(targets []*ent.VersionTarget, processContents func(string, io.Reader) (io.Reader, error)) ([]byte, error) {
	finalZip := bytes.NewBuffer([]byte{})
	zipWriter := zip.NewWriter(finalZip)

	for _, target := range targets {
		archive, err := getTargetArchive(target)
		if err != nil {
			return nil, fmt.Errorf("failed to get target %s: %w", target.TargetName, err)
		}

		zipReader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
		if err != nil {
			return nil, fmt.Errorf("failed to open zip: %w", err)
		}

		for _, file := range zipReader.File {
			err := copyFile(zipWriter, file, target.TargetName+"/", processContents)
			if err != nil {
				return nil, fmt.Errorf("failed to copy file: %w", err)
			}
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close zip: %w", err)
	}

	return finalZip.Bytes(), nil
}

func getTargetArchive(target *ent.VersionTarget) ([]byte, error) {
	downloadURL := target.Key

	resp, err := http.Get(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download target: %w", err)
	}
	defer resp.Body.Close()

	archive, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read target: %w", err)
	}

	return archive, nil
}

func copyFile(zip *zip.Writer, file *zip.File, prefix string, processContents func(string, io.Reader) (io.Reader, error)) error {
	fileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()

	newHeader := file.FileHeader
	newHeader.Name = prefix + newHeader.Name

	fileWriter, err := zip.CreateHeader(&newHeader)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	var finalReader io.Reader = fileReader

	if processContents != nil {
		finalReader, err = processContents(file.Name, fileReader)
		if err != nil {
			return fmt.Errorf("failed to process file: %w", err)
		}
	}

	_, err = io.Copy(fileWriter, finalReader)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func processSMLUplugin(factoryGameVersion string) func(string, io.Reader) (io.Reader, error) {
	return func(name string, file io.Reader) (io.Reader, error) {
		if name != "SML.uplugin" {
			return file, nil
		}

		upluginBytes, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		var uplugin map[string]interface{}
		if err := json.Unmarshal(upluginBytes, &uplugin); err != nil {
			return nil, fmt.Errorf("failed to parse uplugin: %w", err)
		}

		uplugin["GameVersion"] = factoryGameVersion

		if _, ok := uplugin["Plugins"]; ok {
			plugins := uplugin["Plugins"].([]interface{})

			for _, plugin := range plugins {
				pluginData := plugin.(map[string]interface{})
				pluginData["BasePlugin"] = true
			}
		}

		upluginBytes, err = json.Marshal(uplugin)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize uplugin: %w", err)
		}

		return bytes.NewReader(upluginBytes), nil
	}
}
