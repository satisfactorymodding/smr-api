package tests

import (
	"strconv"
	"testing"
	"time"

	"github.com/MarvinJWendt/testza"

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

func TestSMLVersions(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	secondaryVersions := [2]string{
		"2.3.4",
		"3.4.5",
	}

	// Run Twice to detect any cache issues
	for i := 0; i < 2; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var objID string

			t.Run("Create", func(t *testing.T) {
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
					engine_version: "5.1"
				  }) {
					id
				  }
				}`, token)

				var createResponse struct {
					CreateSMLVersion generated.SMLVersion
				}
				testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))

				testza.AssertNotEqual(t, "", createResponse.CreateSMLVersion.ID)

				objID = createResponse.CreateSMLVersion.ID
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
				queryRequest.Var("id", objID)

				var queryResponse struct {
					GetSMLVersion generated.SMLVersion
				}
				testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))

				testza.AssertEqual(t, objID, queryResponse.GetSMLVersion.ID)
				testza.AssertEqual(t, "1.2.3", queryResponse.GetSMLVersion.Version)
				testza.AssertEqual(t, 123, queryResponse.GetSMLVersion.SatisfactoryVersion)
				testza.AssertEqual(t, generated.VersionStabilitiesRelease, queryResponse.GetSMLVersion.Stability)
				testza.AssertEqual(t, "https://google.com", queryResponse.GetSMLVersion.Link)
				testza.AssertEqual(t, "Hello World", queryResponse.GetSMLVersion.Changelog)
				testza.AssertEqual(t, "0.0.0", *queryResponse.GetSMLVersion.BootstrapVersion)
				testza.AssertEqual(t, "5.1", queryResponse.GetSMLVersion.EngineVersion)

				date, err := time.Parse(time.RFC3339, queryResponse.GetSMLVersion.Date)
				testza.AssertNoError(t, err)

				realDate, _ := time.Parse(time.RFC3339, "2023-10-27T01:00:51+00:00")
				testza.AssertEqual(t, realDate.Unix(), date.Unix())

				testza.AssertEqual(t, []*generated.SMLVersionTarget{
					{
						TargetName: generated.TargetNameWindows,
						Link:       "https://this-is-windows.com",
					},
					{
						TargetName: generated.TargetNameWindowsServer,
						Link:       "https://this-is-windows-server.com",
					},
					{
						TargetName: generated.TargetNameLinuxServer,
						Link:       "https://this-is-linux-server.com",
					},
				}, queryResponse.GetSMLVersion.Targets)
			})

			t.Run("Update", func(t *testing.T) {
				updateRequest := authRequest(`mutation ($id: SMLVersionID!, $version: String!) {
				  updateSMLVersion(
					smlVersionId: $id,
					smlVersion: {
					  version: $version,
					  satisfactory_version: 234,
					  stability: alpha,
					  link: "https://ficsit.app",
					  targets: [
						{
						  targetName: Windows,
						  link: "https://this-is-windows-2.com"
						},
						{
						  targetName: WindowsServer,
						  link: "https://this-is-windows-server-2.com"
						},
						{
						  targetName: LinuxServer,
						  link: "https://this-is-linux-server-2.com"
						}
					  ],
					  changelog: "Foo Bar",
					  date: "2000-10-27T01:00:51+00:00",
					  bootstrap_version: "0.0.0",
					  engine_version: "5.2"
						}
				  ) {
					id
				  }
				}`, token)
				updateRequest.Var("id", objID)
				updateRequest.Var("version", secondaryVersions[i])

				var updateResponse struct {
					UpdateSMLVersion generated.SMLVersion
				}
				testza.AssertNoError(t, client.Run(ctx, updateRequest, &updateResponse))
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
				testza.AssertEqual(t, objID, queryResponse.GetSMLVersions.SmlVersions[0].ID)
				testza.AssertEqual(t, secondaryVersions[i], queryResponse.GetSMLVersions.SmlVersions[0].Version)
				testza.AssertEqual(t, 234, queryResponse.GetSMLVersions.SmlVersions[0].SatisfactoryVersion)
				testza.AssertEqual(t, generated.VersionStabilitiesAlpha, queryResponse.GetSMLVersions.SmlVersions[0].Stability)
				testza.AssertEqual(t, "https://ficsit.app", queryResponse.GetSMLVersions.SmlVersions[0].Link)
				testza.AssertEqual(t, "Foo Bar", queryResponse.GetSMLVersions.SmlVersions[0].Changelog)
				testza.AssertEqual(t, "0.0.0", *queryResponse.GetSMLVersions.SmlVersions[0].BootstrapVersion)
				testza.AssertEqual(t, "5.2", queryResponse.GetSMLVersions.SmlVersions[0].EngineVersion)

				date, err := time.Parse(time.RFC3339, queryResponse.GetSMLVersions.SmlVersions[0].Date)
				testza.AssertNoError(t, err)

				realDate, _ := time.Parse(time.RFC3339, "2000-10-27T01:00:51+00:00")
				testza.AssertEqual(t, realDate.Unix(), date.Unix())

				testza.AssertEqual(t, []*generated.SMLVersionTarget{
					{
						TargetName: generated.TargetNameWindows,
						Link:       "https://this-is-windows-2.com",
					},
					{
						TargetName: generated.TargetNameWindowsServer,
						Link:       "https://this-is-windows-server-2.com",
					},
					{
						TargetName: generated.TargetNameLinuxServer,
						Link:       "https://this-is-linux-server-2.com",
					},
				}, queryResponse.GetSMLVersions.SmlVersions[0].Targets)
			})

			t.Run("Delete", func(t *testing.T) {
				deleteRequest := authRequest(`mutation ($id: SMLVersionID!) {
				  deleteSMLVersion(smlVersionId: $id)
				}`, token)
				deleteRequest.Var("id", objID)

				var deleteResponse struct {
					DeleteSMLVersion bool
				}
				testza.AssertNoError(t, client.Run(ctx, deleteRequest, &deleteResponse))

				testza.AssertTrue(t, deleteResponse.DeleteSMLVersion)
			})
		})
	}
}
