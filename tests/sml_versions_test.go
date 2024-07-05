package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/MarvinJWendt/testza"
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"

	_ "embed"
)

func init() {
	config.SetConfigDir("../")
	db.EnableDebug()
}

const smlTestModPath = "testdata/SML-3.7.0.smod"

func TestSMLVersions(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	viper.Set("skip-virus-check", true)

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	t.Run("Create Satisfactory Version", func(t *testing.T) {
		createRequest := authRequest(`mutation {
		  createSatisfactoryVersion(input: {
			version: 123456,
			engine_version: "5.2"
		  }) {
			id
		  }
		}`, token)

		var createResponse struct {
			CreateSatisfactoryVersion generated.SatisfactoryVersion
		}
		testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
		testza.AssertNotEqual(t, "", createResponse.CreateSatisfactoryVersion.ID)
	})

	var smlModID string

	t.Run("Get SML Mod", func(t *testing.T) {
		getRequest := authRequest(`{
   		  getModByReference(modReference: "SML") {
			id
          }
 	    }`, token)

		var getResponse struct {
			GetModByReference generated.Mod
		}
		testza.AssertNoError(t, client.Run(ctx, getRequest, &getResponse))
		smlModID = getResponse.GetModByReference.ID
	})

	var versionID string
	var versionDate time.Time

	t.Run("Create", func(t *testing.T) {
		t.Run("Create Version", func(t *testing.T) {
			createRequest := authRequest(`mutation CreateVersion($mod_id: ModID!) {
				createVersion(modId: $mod_id)
			}`, token)
			createRequest.Var("mod_id", smlModID)

			var createResponse struct {
				CreateVersion string
			}
			testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
			testza.AssertNotEqual(t, "", createResponse.CreateVersion)

			versionID = createResponse.CreateVersion
		})

		t.Run("Upload Parts", func(t *testing.T) {
			f, err := os.Open(smlTestModPath)
			testza.AssertNoError(t, err)

			stat, err := f.Stat()
			testza.AssertNoError(t, err)

			chunkSize := int64(1e+7)
			chunkCount := int(math.Ceil(float64(stat.Size()) / float64(chunkSize))) // Split in 10MB chunks

			for i := range chunkCount {
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
							"mod_id":     smlModID,
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

					part, err := writer.CreateFormFile("0", path.Base(smlTestModPath))
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
			finalizeRequest.Var("mod_id", smlModID)
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
						created_at
					}
					auto_approved
				}
			}`, token)
			request.Var("mod_id", smlModID)
			request.Var("version_id", versionID)

			end := time.Now().Add(time.Minute * 5)
			for time.Now().Before(end) {
				var response struct {
					CheckVersionUploadState struct {
						Version struct {
							ID        string
							CreatedAt string `json:"created_at"`
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
					versionID = response.CheckVersionUploadState.Version.ID
					date, err := time.Parse(time.RFC3339, response.CheckVersionUploadState.Version.CreatedAt)
					testza.AssertNoError(t, err)
					versionDate = date
					break
				}

				time.Sleep(time.Second * 3)
			}

			if time.Now().After(end) {
				testza.AssertNoError(t, errors.New("failed finishing mod"))
			}
		})
	})

	t.Run("Query One", func(t *testing.T) {
		queryRequest := authRequest(`query ($id: SMLVersionID!) {
		  getSMLVersion(smlVersionID: $id) {
			id
			version
			satisfactory_version
			stability
			link
			targets {
			  targetName
			  link
			}
			changelog
			date
			bootstrap_version
			engine_version
		  }
		}`, token)
		queryRequest.Var("id", versionID)

		var queryResponse struct {
			GetSMLVersion generated.SMLVersion
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))

		testza.AssertEqual(t, versionID, queryResponse.GetSMLVersion.ID)
		testza.AssertEqual(t, "3.7.0", queryResponse.GetSMLVersion.Version)
		testza.AssertEqual(t, 273254, queryResponse.GetSMLVersion.SatisfactoryVersion)
		testza.AssertEqual(t, generated.VersionStabilitiesRelease, queryResponse.GetSMLVersion.Stability)
		testza.AssertEqual(t, "https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/tag/v3.7.0", queryResponse.GetSMLVersion.Link)
		testza.AssertEqual(t, "Hello World", queryResponse.GetSMLVersion.Changelog)
		testza.AssertNil(t, queryResponse.GetSMLVersion.BootstrapVersion)
		testza.AssertEqual(t, "5.2", queryResponse.GetSMLVersion.EngineVersion)

		date, err := time.Parse(time.RFC3339, queryResponse.GetSMLVersion.Date)
		testza.AssertNoError(t, err)

		testza.AssertEqual(t, versionDate.Unix(), date.Unix())

		testza.AssertEqual(t, []*generated.SMLVersionTarget{
			{
				TargetName: generated.TargetNameLinuxServer,
				Link:       "https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v3.7.0/SML-LinuxServer.zip",
			},
			{
				TargetName: generated.TargetNameWindows,
				Link:       "https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v3.7.0/SML-Windows.zip",
			},
			{
				TargetName: generated.TargetNameWindowsServer,
				Link:       "https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v3.7.0/SML-WindowsServer.zip",
			},
		}, queryResponse.GetSMLVersion.Targets)
	})

	t.Run("Query Many", func(t *testing.T) {
		queryRequest := authRequest(`{
		  getSMLVersions {
			count
			sml_versions {
			  id
			  version
			  satisfactory_version
			  stability
			  link
			  targets {
				targetName
				link
			  }
			  changelog
			  date
			  bootstrap_version
			  engine_version
			}
		  }
		}`, token)

		var queryResponse struct {
			GetSMLVersions generated.GetSMLVersions
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))

		testza.AssertEqual(t, 1, queryResponse.GetSMLVersions.Count)
		testza.AssertEqual(t, versionID, queryResponse.GetSMLVersions.SmlVersions[0].ID)
		testza.AssertEqual(t, "3.7.0", queryResponse.GetSMLVersions.SmlVersions[0].Version)
		testza.AssertEqual(t, 273254, queryResponse.GetSMLVersions.SmlVersions[0].SatisfactoryVersion)
		testza.AssertEqual(t, generated.VersionStabilitiesRelease, queryResponse.GetSMLVersions.SmlVersions[0].Stability)
		testza.AssertEqual(t, "https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/tag/v3.7.0", queryResponse.GetSMLVersions.SmlVersions[0].Link)
		testza.AssertEqual(t, "Hello World", queryResponse.GetSMLVersions.SmlVersions[0].Changelog)
		testza.AssertNil(t, queryResponse.GetSMLVersions.SmlVersions[0].BootstrapVersion)
		testza.AssertEqual(t, "5.2", queryResponse.GetSMLVersions.SmlVersions[0].EngineVersion)

		date, err := time.Parse(time.RFC3339, queryResponse.GetSMLVersions.SmlVersions[0].Date)
		testza.AssertNoError(t, err)

		testza.AssertEqual(t, versionDate.Unix(), date.Unix())

		testza.AssertEqual(t, []*generated.SMLVersionTarget{
			{
				TargetName: generated.TargetNameLinuxServer,
				Link:       "https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v3.7.0/SML-LinuxServer.zip",
			},
			{
				TargetName: generated.TargetNameWindows,
				Link:       "https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v3.7.0/SML-Windows.zip",
			},
			{
				TargetName: generated.TargetNameWindowsServer,
				Link:       "https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v3.7.0/SML-WindowsServer.zip",
			},
		}, queryResponse.GetSMLVersions.SmlVersions[0].Targets)
	})
}
