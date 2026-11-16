package validation

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/Vilsol/slox"
	"github.com/avast/retry-go"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/proto/parser"
	"github.com/satisfactorymodding/smr-api/storage"
)

var AllowedTargets = []string{"Windows", "WindowsServer", "LinuxServer"}

type ModType int

const (
	DataJSON            ModType = iota
	UEPlugin                    = 1
	MultiTargetUEPlugin         = 2
)

type ModMetadata []map[string]map[string][]interface{}

type ModInfo struct {
	Dependencies         map[string]string `json:"dependencies"`
	OptionalDependencies map[string]string `json:"optional_dependencies"`
	Semver               *semver.Version   `json:"semver"`
	ModReference         string            `json:"mod_reference"`
	Version              string            `json:"version"`
	Hash                 string            `json:"hash"`
	GameVersion          string            `json:"game_version"`
	RequiredOnRemote     bool              `json:"required_on_remote"`
	Targets              []string          `json:"targets"`
	FilesToScan          []string          `json:"files_to_scan"`
	Size                 int64             `json:"size"`
	Type                 ModType           `json:"type"`
}

var (
	dataJSONSchema    gojsonschema.JSONLoader
	uPluginJSONSchema gojsonschema.JSONLoader
)

var StaticPath = "static/"

func InitializeValidator() {
	absPath, err := filepath.Abs(filepath.Join(StaticPath, "data-json-schema.json"))
	if err != nil {
		panic(err)
	}

	dataJSONSchema = gojsonschema.NewReferenceLoader("file://" + strings.ReplaceAll(absPath, "\\", "/"))

	absPath, err = filepath.Abs(filepath.Join(StaticPath, "uplugin-json-schema.json"))
	if err != nil {
		panic(err)
	}

	uPluginJSONSchema = gojsonschema.NewReferenceLoader("file://" + strings.ReplaceAll(absPath, "\\", "/"))
}

func ExtractModInfo(ctx context.Context, body []byte, withValidation bool, modReference string) (*ModInfo, error) {
	ctx, span := otel.Tracer("ficsit-app").Start(ctx, "ExtractModInfo")
	defer span.End()

	if len(body) > 1000000000 {
		err := errors.New("mod archive must be < 1GB")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	archive, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, errors.New("invalid zip archive")
	}

	var dataFile *zip.File
	var uPlugin *zip.File

	for _, v := range archive.File {
		if v.Name == "data.json" {
			dataFile = v
			break
		}
		if v.Name == modReference+".uplugin" {
			uPlugin = v
			break
		}
	}

	var modInfo *ModInfo

	if dataFile != nil {
		modInfo, err = validateDataJSON(ctx, archive, dataFile, withValidation)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return nil, err
		}
	}

	if uPlugin != nil {
		modInfo, err = validateUPluginJSON(ctx, uPlugin, withValidation, modReference)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return nil, err
		}
	}

	if modInfo == nil {
		// Neither data.json nor .uplugin found, try multi-target .uplugin
		modInfo, err = validateMultiTargetPlugin(ctx, archive, withValidation, modReference)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return nil, err
		}
	}

	if modInfo == nil {
		err := errors.New("missing " + modReference + ".uplugin or data.json")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	modInfo.Size = int64(len(body))

	hash := sha256.New()
	_, err = hash.Write(body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		slox.Error(ctx, "error hashing pak", slog.Any("err", err))
	}

	modInfo.Hash = hex.EncodeToString(hash.Sum(nil))

	version, err := semver.StrictNewVersion(modInfo.Version)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		slox.Error(ctx, "error parsing semver", slog.Any("err", err))
		return nil, fmt.Errorf("error parsing semver: %w", err)
	}

	modInfo.Semver = version

	modInfo.FilesToScan, err = getFilesToScan(archive)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		slox.Error(ctx, "error getting files to scan", slog.Any("err", err))
		return nil, fmt.Errorf("error getting files to scan: %w", err)
	}

	return modInfo, nil
}

func ExtractMetadata(ctx context.Context, data []byte, gameVersion string, modReference string) (ModMetadata, error) {
	ctx, span := otel.Tracer("ficsit-app").Start(ctx, "extractMetadata")
	defer span.End()

	metadata := make(ModMetadata, 0)

	// Extract all possible metadata
	conn, err := grpc.NewClient(
		viper.GetString("extractor_host"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, fmt.Errorf("failed to connect to metadata server: %w", err)
	}
	defer conn.Close()

	engineVersion := "4.26"

	//nolint
	if db.From(ctx) != nil {
		engineVersion, err = db.GetEngineVersionForSatisfactoryVersion(ctx, gameVersion)
		if err != nil {
			slox.Warn(ctx, "failed to get engine version", slog.Any("err", err))
		}
	} else {
		slox.Warn(ctx, "no database context provided to validator")
	}

	slox.Info(ctx, "decided engine version", slog.String("version", engineVersion))

	if err := retry.Do(func() error {
		parserClient := parser.NewParserClient(conn)
		stream, err := parserClient.Parse(ctx, &parser.ParseRequest{
			ZipData:       data,
			EngineVersion: engineVersion,
		},
			grpc.MaxCallSendMsgSize(1024*1024*1024), // 1GB
			grpc.MaxCallRecvMsgSize(1024*1024*1024), // 1GB
		)
		if err != nil {
			return fmt.Errorf("failed to parse mod: %w", err)
		}

		defer func(stream parser.Parser_ParseClient) {
			err := stream.CloseSend()
			if err != nil {
				slox.Error(ctx, "failed closing parser stream", slog.Any("err", err))
			}
		}(stream)

		beforeUpload := time.Now().Add(-time.Minute)

		count := 0
		for {
			asset, err := stream.Recv()
			if err != nil {
				//nolint
				if errors.Is(err, io.EOF) || err == io.EOF {
					break
				}
				return fmt.Errorf("failed reading parser stream: %w", err)
			}

			slox.Info(ctx, "received asset from parser", slog.String("path", asset.GetPath()))

			if asset.Path == "metadata.json" {
				out, err := ExtractMetadataRaw(asset.Data)
				if err != nil {
					return err
				}
				metadata = append(metadata, out)
			}

			storage.UploadModAsset(ctx, modReference, asset.GetPath(), asset.GetData())
			count++
		}

		slox.Info(ctx, "all assets received", slog.Int("count", count))

		storage.DeleteOldModAssets(ctx, modReference, beforeUpload)

		return nil
	},
		retry.Attempts(10),
		retry.Delay(time.Second*10),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			if n > 0 {
				slox.Info(ctx, "retrying to extract metadata", slog.Uint64("n", uint64(n)), slog.Any("err", err))
			}
		})); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err //nolint
	}

	return metadata, nil
}

type DataJSONInfo struct {
	Dependencies         map[string]string   `json:"dependencies"`
	OptionalDependencies map[string]string   `json:"optional_dependencies"`
	ModReference         string              `json:"mod_reference"`
	Version              string              `json:"version"`
	SMLVersion           *string             `json:"sml_version"`
	Objects              []DataJSONModObject `json:"objects"`
}

type DataJSONModObject struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

func validateDataJSON(ctx context.Context, archive *zip.Reader, dataFile *zip.File, withValidation bool) (*ModInfo, error) {
	rc, err := dataFile.Open()
	defer func(rc io.ReadCloser) {
		_ = rc.Close()
	}(rc)

	if err != nil {
		return nil, errors.New("invalid zip archive")
	}

	dataJSON, err := io.ReadAll(rc)
	if err != nil {
		return nil, errors.New("invalid zip archive")
	}

	result, err := gojsonschema.Validate(dataJSONSchema, gojsonschema.NewBytesLoader(dataJSON))
	if err != nil {
		return nil, errors.New("data.json doesn't follow schema. please view the help page. (" + err.Error() + ")")
	}

	if withValidation {
		if !result.Valid() {
			return nil, errors.New("data.json doesn't follow schema. please view the help page. (" + fmt.Sprintf("%s", result.Errors()) + ")")
		}
	}

	var dataJSONInfo DataJSONInfo
	err = json.Unmarshal(dataJSON, &dataJSONInfo)
	if err != nil {
		return nil, errors.New("invalid data.json")
	}

	if withValidation {
		if len(dataJSONInfo.Dependencies) == 0 {
			return nil, errors.New("data.json doesn't contain SML as a dependency.") //nolint:revive
		}
	}

	if smlDep, ok := dataJSONInfo.Dependencies["SML"]; ok {
		dataJSONInfo.SMLVersion = &smlDep
	}

	if dataJSONInfo.SMLVersion == nil {
		return nil, errors.New("data.json doesn't contain SML as a dependency.") //nolint:revive
	}

	// Validate that all listed files are accounted for in data.json
	for _, archiveFile := range archive.File {
		if archiveFile != nil {
			if strings.HasSuffix(archiveFile.Name, ".dll") || strings.HasSuffix(archiveFile.Name, ".pak") || strings.HasSuffix(archiveFile.Name, ".so") {
				found := false
				for _, obj := range dataJSONInfo.Objects {
					if obj.Path == archiveFile.Name {
						found = true
						break
					}
				}
				if !found {
					return nil, errors.New("zip archive contains unreferenced objects: " + archiveFile.Name)
				}
			}
		}
	}

	// Validate that all objects refer to existing files
	for _, obj := range dataJSONInfo.Objects {
		found := false
		for _, archiveFile := range archive.File {
			if obj.Path == archiveFile.Name {
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("data.json objects refer to non-existent path: " + obj.Path)
		}
	}

	modInfo := ModInfo{
		Dependencies:         dataJSONInfo.Dependencies,
		OptionalDependencies: dataJSONInfo.OptionalDependencies,
		ModReference:         dataJSONInfo.ModReference,
		Version:              dataJSONInfo.Version,
		RequiredOnRemote:     true,
		Targets:              []string{"Windows"},
		Type:                 DataJSON,
	}

	gameVersion, err := getGameVersionFromSMLVersion(ctx, dataJSONInfo.Dependencies["SML"])
	if err != nil {
		return nil, fmt.Errorf("failed to infer FactoryGame version: %w", err)
	}
	modInfo.GameVersion = gameVersion

	return &modInfo, nil
}

type UPlugin struct {
	SemVersion       *string  `json:"SemVersion"`
	Plugins          []Plugin `json:"Plugins"`
	Version          int64    `json:"Version"`
	GameVersion      string   `json:"GameVersion"`
	RequiredOnRemote *bool    `json:"RequiredOnRemote"`
}

type Plugin struct {
	BasePlugin *bool  `json:"BasePlugin"`
	Optional   *bool  `json:"Optional"`
	Name       string `json:"Name"`
	SemVersion string `json:"SemVersion"`
}

func validateUPluginJSON(ctx context.Context, uPluginFile *zip.File, withValidation bool, modReference string) (*ModInfo, error) {
	rc, err := uPluginFile.Open()
	defer func(rc io.ReadCloser) {
		_ = rc.Close()
	}(rc)

	if err != nil {
		return nil, errors.New("invalid zip archive")
	}

	uPluginJSON, err := io.ReadAll(rc)
	if err != nil {
		return nil, errors.New("invalid zip archive")
	}

	result, err := gojsonschema.Validate(uPluginJSONSchema, gojsonschema.NewBytesLoader(uPluginJSON))
	if err != nil {
		return nil, errors.New(uPluginFile.Name + " doesn't follow schema. please view the help page. (" + err.Error() + ")")
	}

	if withValidation {
		if !result.Valid() {
			return nil, errors.New(uPluginFile.Name + " doesn't follow schema. please view the help page. (" + fmt.Sprintf("%s", result.Errors()) + ")")
		}
	}

	var uPlugin UPlugin
	err = json.Unmarshal(uPluginJSON, &uPlugin)
	if err != nil {
		return nil, errors.New("invalid " + uPluginFile.Name)
	}

	modInfo := ModInfo{
		ModReference:         modReference,
		Dependencies:         map[string]string{},
		OptionalDependencies: map[string]string{},
		RequiredOnRemote:     true,
		Type:                 UEPlugin,
	}

	if uPlugin.SemVersion != nil {
		modInfo.Version = *uPlugin.SemVersion

		split := strings.Split(modInfo.Version, ".")
		if split[0] != strconv.FormatInt(uPlugin.Version, 10) {
			return nil, errors.New("SemVer major version should match Version")
		}
	} else {
		modInfo.Version = strconv.FormatInt(uPlugin.Version, 10) + ".0.0"
	}

	if uPlugin.RequiredOnRemote != nil {
		modInfo.RequiredOnRemote = *uPlugin.RequiredOnRemote
	}

	for _, plugin := range uPlugin.Plugins {
		if plugin.BasePlugin != nil && *plugin.BasePlugin {
			continue
		}

		if plugin.Optional != nil && *plugin.Optional {
			modInfo.OptionalDependencies[plugin.Name] = plugin.SemVersion
		} else {
			modInfo.Dependencies[plugin.Name] = plugin.SemVersion
		}
	}

	if uPlugin.GameVersion != "" {
		modInfo.GameVersion = uPlugin.GameVersion
	} else if smlVersion, ok := modInfo.Dependencies["SML"]; ok {
		gameVersion, err := getGameVersionFromSMLVersion(ctx, smlVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to infer FactoryGame version: %w", err)
		}
		modInfo.GameVersion = gameVersion
	} else {
		return nil, fmt.Errorf("infering FactoryGame version: %s doesn't contain SML as a dependency", uPluginFile.Name)
	}

	return &modInfo, nil
}

func getGameVersionFromSMLVersion(ctx context.Context, smlVersion string) (string, error) {
	smlQuery := db.From(ctx).Mod.Query().Where(mod.ModReferenceEQ("SML")).WithVersions()
	smlMod, err := smlQuery.First(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get SML mod: %w", err)
	}

	smlVersions := smlMod.Edges.Versions

	// Sort increasing by version
	sort.Slice(smlVersions, func(a, b int) bool {
		return semver.MustParse(smlVersions[a].Version).Compare(semver.MustParse(smlVersions[b].Version)) < 0
	})

	constraint, err := semver.NewConstraint(smlVersion)
	if err != nil {
		return "", fmt.Errorf("failed to create semver constraint: %w", err)
	}

	for _, version := range smlVersions {
		if constraint.Check(semver.MustParse(version.Version)) {
			return version.GameVersion, nil
		}
	}

	return "", fmt.Errorf("no SML version matches constraint: %s", smlVersion)
}

func validateMultiTargetPlugin(ctx context.Context, archive *zip.Reader, withValidation bool, modReference string) (*ModInfo, error) {
	var targets []string
	var uPluginFiles []*zip.File
	for _, file := range archive.File {
		if path.Base(file.Name) == modReference+".uplugin" && path.Dir(file.Name) != "." {
			targets = append(targets, path.Dir(file.Name))
			uPluginFiles = append(uPluginFiles, file)
		}
	}

	if withValidation {
		for _, target := range targets {
			found := false
			for _, allowedTarget := range AllowedTargets {
				if target == allowedTarget {
					found = true
					break
				}
			}
			if !found {
				return nil, errors.New("multi-target plugin contains invalid target: " + target)
			}
		}

		for _, file := range archive.File {
			found := false
			for _, target := range targets {
				if strings.HasPrefix(file.Name, target+"/") {
					found = true
					break
				}
			}
			if !found {
				return nil, errors.New("multi-target plugin contains file outside of target directories: " + file.Name)
			}
		}
	}

	if len(uPluginFiles) == 0 {
		return nil, errors.New("multi-target plugin doesn't contain any .uplugin files")
	}

	if withValidation {
		var lastData []byte
		for _, uPluginFile := range uPluginFiles {
			file, err := uPluginFile.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open .uplugin file: %w", err)
			}
			data, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read .uplugin file: %w", err)
			}

			if lastData != nil && !bytes.Equal(lastData, data) {
				return nil, errors.New("multi-target plugin contains different .uplugin files")
			}
			lastData = data
		}
	}

	// All the .uplugin files should be the same at this point (assuming validation is enabled)
	modInfo, err := validateUPluginJSON(ctx, uPluginFiles[0], withValidation, modReference)
	if err != nil {
		return nil, fmt.Errorf("failed to validate multi-target plugin: %w", err)
	}

	modInfo.Targets = targets
	modInfo.Type = MultiTargetUEPlugin

	return modInfo, nil
}

func getFilesToScan(archive *zip.Reader) ([]string, error) {
	filesToScan := make([]string, 0)

	for _, file := range archive.File {
		if file.Mode().IsDir() {
			continue
		}

		scan, err := shouldScanFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to determine if file should be scanned: %w", err)
		}

		if scan {
			filesToScan = append(filesToScan, file.Name)
		}
	}

	return filesToScan, nil
}

var scanMimeTypes = []string{
	"application/vnd.microsoft.portable-executable",
	"application/x-elf", // shared objects and other elf types are children of this mime type
}

func shouldScanFile(file *zip.File) (bool, error) {
	f, err := file.Open()
	if err != nil {
		return false, fmt.Errorf("failed to open file %s: %w", file.Name, err)
	}
	defer f.Close()

	mime, err := mimetype.DetectReader(f)
	if err != nil {
		return false, fmt.Errorf("failed to detect mime type for file %s: %w", file.Name, err)
	}

	for mime != nil {
		if mimetype.EqualsAny(mime.String(), scanMimeTypes...) {
			return true, nil
		}
		mime = mime.Parent()
	}

	return false, nil
}
