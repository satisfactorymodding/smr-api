package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/MarvinJWendt/testza"
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/migrations"
)

func init() {
	migrations.SetMigrationDir("../migrations")
	config.SetConfigDir("../")
	postgres.EnableDebug()
	db.EnableDebug()
}

const testModURL = "https://storage-staging.ficsit.app/file/smr-staging-s3/mods/4ha9SfdXaHAq4q/FicsitRemoteMonitoring-0.10.3.smod"

// TODO Add rate limit test

func TestVersions(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	viper.Set("skip-virus-check", true)

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	var modID string

	modReference := "FicsitRemoteMonitoring"

	u, err := url.Parse(testModURL)
	testza.AssertNoError(t, err)

	modPath := filepath.Join(t.TempDir(), path.Base(u.Path))

	t.Run("Download Test Mod", func(t *testing.T) {
		response, err := http.Get(testModURL)
		testza.AssertNoError(t, err)

		defer response.Body.Close()

		file, err := os.OpenFile(modPath, os.O_CREATE|os.O_WRONLY, 0o644)
		testza.AssertNoError(t, err)

		_, err = io.Copy(file, response.Body)
		testza.AssertNoError(t, err)
	})

	t.Run("Create SML Version", func(t *testing.T) {
		createRequest := authRequest(`mutation {
		  createSMLVersion(smlVersion: {
			version: "1.2.3",
			satisfactory_version: 123,
			stability: release,
			link: "https://google.com",
			targets: [
			  {
				targetName: Windows,
				link: "https://this-is-windows.com"
			  },
			  {
				targetName: WindowsServer,
				link: "https://this-is-windows-server.com"
			  },
			  {
				targetName: LinuxServer,
				link: "https://this-is-linux-server.com"
			  }
			],
			changelog: "Hello World",
			date: "2023-10-27T01:00:51+00:00",
			bootstrap_version: "0.0.0",
			engine_version: "5.2"
		  }) {
			id
		  }
		}`, token)

		var createResponse struct {
			CreateSMLVersion generated.SMLVersion
		}
		testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
		testza.AssertNotEqual(t, "", createResponse.CreateSMLVersion.ID)
	})

	t.Run("Create Mod", func(t *testing.T) {
		createRequest := authRequest(`mutation CreateMod($mod_reference: ModReference!) {
			createMod(mod: {
				name: "Hello World",
				short_description: "Foo Bar 123 Foo Bar 123",
				full_description: "Lorem ipsum dolor sit amet",
				mod_reference: $mod_reference
			}) {
				id
			}
		}`, token)
		createRequest.Var("mod_reference", modReference)

		var createResponse struct {
			CreateMod generated.Mod
		}
		testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
		testza.AssertNotEqual(t, "", createResponse.CreateMod.ID)

		modID = createResponse.CreateMod.ID
	})

	var versionID string

	t.Run("Create Version", func(t *testing.T) {
		createRequest := authRequest(`mutation CreateVersion($mod_id: ModID!) {
			createVersion(modId: $mod_id)
		}`, token)
		createRequest.Var("mod_id", modID)

		var createResponse struct {
			CreateVersion string
		}
		testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
		testza.AssertNotEqual(t, "", createResponse.CreateVersion)

		versionID = createResponse.CreateVersion
	})

	t.Run("Upload Parts", func(t *testing.T) {
		f, err := os.Open(modPath)
		testza.AssertNoError(t, err)

		stat, err := f.Stat()
		testza.AssertNoError(t, err)

		chunkSize := int64(1e+7)
		chunkCount := int(math.Ceil(float64(stat.Size()) / float64(chunkSize))) // Split in 10MB chunks

		for i := 0; i < chunkCount; i++ {
			t.Run("Part"+strconv.Itoa(i), func(t *testing.T) {
				_, err = f.Seek(int64(i)*chunkSize, 0)
				testza.AssertNoError(t, err)

				chunk := make([]byte, chunkSize)
				n, err := f.Read(chunk)
				testza.AssertNoError(t, err)
				chunk = chunk[:n]

				operationBody, err := json.Marshal(map[string]interface{}{
					"query": `mutation UploadVersionPart($mod_id: ModID!, $version_id: VersionID!, $part: Int!, $file: Upload!) {
						uploadVersionPart(
							modId: $mod_id,
							versionId: $version_id,
							file: $file,
							part: $part
						)
					}`,
					"variables": map[string]interface{}{
						"mod_id":     modID,
						"version_id": versionID,
						"part":       i + 1,
						"file":       nil,
					},
				})
				testza.AssertNoError(t, err)

				mapBody, err := json.Marshal(map[string]interface{}{
					"0": []string{"variables.file"},
				})
				testza.AssertNoError(t, err)

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				operations, err := writer.CreateFormField("operations")
				testza.AssertNoError(t, err)

				_, err = operations.Write(operationBody)
				testza.AssertNoError(t, err)

				mapField, err := writer.CreateFormField("map")
				testza.AssertNoError(t, err)

				_, err = mapField.Write(mapBody)
				testza.AssertNoError(t, err)

				part, err := writer.CreateFormFile("0", filepath.Base(path.Base(u.Path)))
				testza.AssertNoError(t, err)

				_, err = io.Copy(part, bytes.NewReader(chunk))
				testza.AssertNoError(t, err)

				err = writer.Close()
				testza.AssertNoError(t, err)

				r, _ := http.NewRequest("POST", "http://localhost:5020/v2/query", body)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				r.Header.Add("Authorization", token)

				resp, err := http.DefaultClient.Do(r)
				testza.AssertNoError(t, err)

				defer resp.Body.Close()
				all, err := io.ReadAll(resp.Body)
				testza.AssertNoError(t, err)

				response := make(map[string]interface{})
				testza.AssertNoError(t, json.Unmarshal(all, &response))

				testza.AssertTrue(t, response["data"].(map[string]interface{})["uploadVersionPart"].(bool))
			})
		}
	})

	t.Run("Finalize Version", func(t *testing.T) {
		finalizeRequest := authRequest(`mutation FinalizeCreateVersion($mod_id: ModID!, $version_id: VersionID!) {
			finalizeCreateVersion(modId: $mod_id, versionId: $version_id, version: {
				changelog: "Hello World",
				stability: release
			})
		}`, token)
		finalizeRequest.Var("mod_id", modID)
		finalizeRequest.Var("version_id", versionID)

		var finalizeResponse struct {
			FinalizeCreateVersion bool
		}
		testza.AssertNoError(t, client.Run(ctx, finalizeRequest, &finalizeResponse))
		testza.AssertTrue(t, finalizeResponse.FinalizeCreateVersion)
	})

	t.Run("Wait For Version", func(t *testing.T) {
		request := authRequest(`query CheckVersionUploadState($mod_id: ModID!, $version_id: VersionID!) {
			checkVersionUploadState(modId: $mod_id, versionId: $version_id) {
				version {
					id
				}
				auto_approved
			}
		}`, token)
		request.Var("mod_id", modID)
		request.Var("version_id", versionID)

		end := time.Now().Add(time.Minute * 5)
		for time.Now().Before(end) {
			var response struct {
				CheckVersionUploadState struct {
					Version struct {
						ID string
					}
					AutoApproved bool
				}
			}

			err := client.Run(ctx, request, &response)
			testza.AssertNoError(t, err)

			if err != nil {
				break
			}

			if response.CheckVersionUploadState.Version.ID != "" {
				break
			}

			time.Sleep(time.Second * 3)
		}

		if time.Now().After(end) {
			testza.AssertNoError(t, errors.New("failed finishing mod"))
		}
	})
}
