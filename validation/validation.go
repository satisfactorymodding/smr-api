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
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/proto/parser"
	"github.com/satisfactorymodding/smr-api/storage"
)

var AllowedTargets = []string{"Windows", "WindowsServer", "LinuxServer"}

type ModObject struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type ModType int

const (
	DataJSON            ModType = iota
	UEPlugin                    = 1
	MultiTargetUEPlugin         = 2
)

type ModInfo struct {
	Dependencies         map[string]string                     `json:"dependencies"`
	OptionalDependencies map[string]string                     `json:"optional_dependencies"`
	Semver               *semver.Version                       `json:"-"`
	ModReference         string                                `json:"mod_reference"`
	Version              string                                `json:"version"`
	Hash                 string                                `json:"-"`
	SMLVersion           string                                `json:"sml_version"`
	Objects              []ModObject                           `json:"objects"`
	Metadata             []map[string]map[string][]interface{} `json:"-"`
	Targets              []string                              `json:"-"`
	Size                 int64                                 `json:"-"`
	Type                 ModType                               `json:"-"`
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

func ExtractModInfo(ctx context.Context, body []byte, withMetadata bool, withValidation bool, modReference string) (*ModInfo, error) {
	if len(body) > 1000000000 {
		return nil, errors.New("mod archive must be < 1GB")
	}

	archive, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
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
		modInfo, err = validateDataJSON(archive, dataFile, withValidation)
		if err != nil {
			return nil, err
		}
	}

	if uPlugin != nil {
		modInfo, err = validateUPluginJSON(archive, uPlugin, withValidation, modReference)
		if err != nil {
			return nil, err
		}
	}

	if modInfo == nil {
		// Neither data.json nor .uplugin found, try multi-target .uplugin
		modInfo, err = validateMultiTargetPlugin(archive, withValidation, modReference)
		if err != nil {
			return nil, err
		}
	}

	if modInfo == nil {
		return nil, errors.New("missing " + modReference + ".uplugin or data.json")
	}

	if withMetadata {
		// Extract all possible metadata
		conn, err := grpc.Dial(viper.GetString("extractor_host"), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to metadata server: %w", err)
		}
		defer conn.Close()

		engineVersion := "4.26"

		//nolint
		if db.From(ctx) != nil {
			smlVersions, err := db.From(ctx).SmlVersion.Query().All(ctx)
			if err != nil {
				return nil, err
			}

			// Sort decrementing by version
			sort.Slice(smlVersions, func(a, b int) bool {
				return semver.MustParse(smlVersions[a].Version).Compare(semver.MustParse(smlVersions[b].Version)) > 0
			})

			for _, version := range smlVersions {
				constraint, err := semver.NewConstraint(modInfo.SMLVersion)
				if err != nil {
					return nil, fmt.Errorf("failed to create semver constraint: %w", err)
				}

				if constraint.Check(semver.MustParse(version.Version)) {
					engineVersion = version.EngineVersion
					break
				}
			}
		}

		slox.Info(ctx, "decided engine version", slog.String("version", engineVersion))

		if err := retry.Do(func() error {
			parserClient := parser.NewParserClient(conn)
			stream, err := parserClient.Parse(ctx, &parser.ParseRequest{
				ZipData:       body,
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
					out, err := ExtractMetadata(asset.Data)
					if err != nil {
						return err
					}
					modInfo.Metadata = append(modInfo.Metadata, out)
				}

				storage.UploadModAsset(ctx, modInfo.ModReference, asset.GetPath(), asset.GetData())
			}

			storage.DeleteOldModAssets(ctx, modInfo.ModReference, beforeUpload)

			return nil
		},
			retry.Attempts(10),
			retry.Delay(time.Second*10),
			retry.DelayType(retry.FixedDelay),
			retry.OnRetry(func(n uint, err error) {
				if n > 0 {
					slox.Info(ctx, "retrying to extract metadata", slog.Uint64("n", uint64(n)))
				}
			})); err != nil {
			return nil, err //nolint
		}
	}

	modInfo.Size = int64(len(body))

	hash := sha256.New()
	_, err = hash.Write(body)

	if err != nil {
		slox.Error(ctx, "error hashing pak", slog.Any("err", err))
	}

	modInfo.Hash = hex.EncodeToString(hash.Sum(nil))

	version, err := semver.StrictNewVersion(modInfo.Version)
	if err != nil {
		slox.Error(ctx, "error parsing semver", slog.Any("err", err))
		return nil, fmt.Errorf("error parsing semver: %w", err)
	}

	modInfo.Semver = version

	return modInfo, nil
}

func validateDataJSON(archive *zip.Reader, dataFile *zip.File, withValidation bool) (*ModInfo, error) {
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

	var modInfo ModInfo
	err = json.Unmarshal(dataJSON, &modInfo)

	if err != nil {
		return nil, errors.New("invalid data.json")
	}

	if withValidation {
		if len(modInfo.Dependencies) == 0 {
			return nil, errors.New("data.json doesn't contain SML as a dependency.") //nolint:revive
		}
	}

	for key, val := range modInfo.Dependencies {
		if key == "SML" {
			modInfo.SMLVersion = val
		}
	}

	if modInfo.SMLVersion == "" {
		return nil, errors.New("data.json doesn't contain SML as a dependency.") //nolint:revive
	}

	// Validate that all listed files are accounted for in data.json
	for _, archiveFile := range archive.File {
		if archiveFile != nil {
			if strings.HasSuffix(archiveFile.Name, ".dll") || strings.HasSuffix(archiveFile.Name, ".pak") || strings.HasSuffix(archiveFile.Name, ".so") {
				found := false
				for _, obj := range modInfo.Objects {
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
	for _, obj := range modInfo.Objects {
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

	modInfo.Type = DataJSON

	return &modInfo, nil
}

type UPlugin struct {
	SemVersion *string  `json:"SemVersion"`
	Plugins    []Plugin `json:"Plugins"`
	Version    int64    `json:"Version"`
}

type Plugin struct {
	BIsBasePlugin *bool  `json:"bIsBasePlugin"`
	BIsOptional   *bool  `json:"bIsOptional"`
	Name          string `json:"Name"`
	SemVersion    string `json:"SemVersion"`
}

func validateUPluginJSON(archive *zip.Reader, uPluginFile *zip.File, withValidation bool, modReference string) (*ModInfo, error) {
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
		Objects:              []ModObject{},
		Dependencies:         map[string]string{},
		OptionalDependencies: map[string]string{},
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

	for _, plugin := range uPlugin.Plugins {
		if plugin.BIsBasePlugin != nil && *plugin.BIsBasePlugin {
			continue
		}

		if plugin.BIsOptional != nil && *plugin.BIsOptional {
			modInfo.OptionalDependencies[plugin.Name] = plugin.SemVersion
		} else {
			modInfo.Dependencies[plugin.Name] = plugin.SemVersion
		}
	}

	for _, file := range archive.File {
		if file != nil {
			splitName := strings.Split(file.Name, ".")
			extension := splitName[len(splitName)-1]
			if extension == "pak" {
				modInfo.Objects = append(modInfo.Objects, ModObject{
					Path: file.Name,
					Type: "pak",
				})
			} else if extension == "dll" || extension == "so" {
				modInfo.Objects = append(modInfo.Objects, ModObject{
					Path: file.Name,
					Type: "sml_mod",
				})
			}
		}
	}

	if withValidation {
		if len(modInfo.Dependencies) == 0 {
			return nil, errors.New(uPluginFile.Name + " doesn't contain SML as a dependency.")
		}
	}

	for key, val := range modInfo.Dependencies {
		if key == "SML" {
			modInfo.SMLVersion = val
		}
	}

	if modInfo.SMLVersion == "" {
		return nil, errors.New(uPluginFile.Name + " doesn't contain SML as a dependency.")
	}

	modInfo.Type = UEPlugin

	return &modInfo, nil
}

func validateMultiTargetPlugin(archive *zip.Reader, withValidation bool, modReference string) (*ModInfo, error) {
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
	modInfo, err := validateUPluginJSON(archive, uPluginFiles[0], withValidation, modReference)
	if err != nil {
		return nil, fmt.Errorf("failed to validate multi-target plugin: %w", err)
	}

	modInfo.Targets = targets
	modInfo.Type = MultiTargetUEPlugin

	return modInfo, nil
}
