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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/Vilsol/ue4pak/parser"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

type ModObject struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type ModInfo struct {
	ModReference         string                                `json:"mod_reference"`
	Version              string                                `json:"version"`
	Objects              []ModObject                           `json:"objects"`
	Dependencies         map[string]string                     `json:"dependencies"`
	OptionalDependencies map[string]string                     `json:"optional_dependencies"`
	Metadata             []map[string]map[string][]interface{} `json:"-"`
	Size                 int64                                 `json:"-"`
	Hash                 string                                `json:"-"`
	Semver               *semver.Version                       `json:"-"`
	SMLVersion           string                                `json:"sml_version"`
}

var dataJSONSchema gojsonschema.JSONLoader
var uPluginJSONSchema gojsonschema.JSONLoader

func InitializeValidator() {
	absPath, err := filepath.Abs("static/data-json-schema.json")

	if err != nil {
		panic(err)
	}

	dataJSONSchema = gojsonschema.NewReferenceLoader("file://" + strings.Replace(absPath, "\\", "/", -1))

	absPath, err = filepath.Abs("static/uplugin-json-schema.json")

	if err != nil {
		panic(err)
	}

	uPluginJSONSchema = gojsonschema.NewReferenceLoader("file://" + strings.Replace(absPath, "\\", "/", -1))
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
		if v.Name == "WindowsNoEditor/"+modReference+".uplugin" {
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
		return nil, errors.New("missing WindowsNoEditor/" + modReference + ".uplugin or data.json")
	}

	if withMetadata {
		// Extract all possible metadata
		modInfo.Metadata = make([]map[string]map[string][]interface{}, 0)
		for _, obj := range modInfo.Objects {
			if strings.ToLower(obj.Type) == "pak" {
				for _, archiveFile := range archive.File {
					if obj.Path == archiveFile.Name {
						data, err := archiveFile.Open()

						if err != nil {
							log.Err(err).Msg("failed opening archive file")
							break
						}

						pakData, err := io.ReadAll(data)

						if err != nil {
							log.Err(err).Msg("failed reading archive file")
							break
						}

						reader := &parser.PakByteReader{
							Bytes: pakData,
						}

						pak, err := AttemptExtractDataFromPak(ctx, reader)

						if err != nil {
							log.Err(err).Msg("failed parsing archive file")
							break
						}

						modInfo.Metadata = append(modInfo.Metadata, pak)
						break
					}
				}
			}
		}
	}

	modInfo.Size = int64(len(body))

	hash := sha256.New()
	_, err = hash.Write(body)

	if err != nil {
		log.Err(err).Msg("error hashing pak")
	}

	modInfo.Hash = hex.EncodeToString(hash.Sum(nil))

	version, err := semver.StrictNewVersion(modInfo.Version)

	if err != nil {
		log.Err(err).Msg("error parsing semver")
		return nil, errors.Wrap(err, "error parsing semver")
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

	return &modInfo, nil
}

type UPlugin struct {
	SemVersion *string  `json:"SemVersion"`
	Version    int64    `json:"Version"`
	Plugins    []Plugin `json:"Plugins"`
}

type Plugin struct {
	Name          string `json:"Name"`
	SemVersion    string `json:"SemVersion"`
	BIsBasePlugin *bool  `json:"bIsBasePlugin"`
	BIsOptional   *bool  `json:"bIsOptional"`
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

	return &modInfo, nil
}
