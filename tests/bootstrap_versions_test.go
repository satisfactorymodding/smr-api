package tests

import (
	"strconv"
	"testing"

	"github.com/MarvinJWendt/testza"

	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/migrations"
)

func init() {
	migrations.SetMigrationDir("../migrations")
	config.SetConfigDir("../")
	postgres.EnableDebug()
}

func TestBootstrapVersions(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	// Run Twice to detect any cache issues
	for i := 0; i < 2; i++ {
		version := strconv.Itoa(i+1) + ".0.0"

		// Create
		createBootstrapVersion := authRequest(`mutation ($version: String!) {
			createBootstrapVersion(bootstrapVersion: {
				version: $version,
				satisfactory_version: 12345,
				stability: beta,
				link: "example.com",
				changelog: "Hello World",
				date: "2006-01-02T15:04:05Z"
			}) {
				id
			}
		}`, token)
		createBootstrapVersion.Var("version", version)

		var createBootstrapVersionResponse struct {
			CreateBootstrapVersion generated.BootstrapVersion
		}
		testza.AssertNoError(t, client.Run(ctx, createBootstrapVersion, &createBootstrapVersionResponse))
		testza.AssertNotEqual(t, "", createBootstrapVersionResponse.CreateBootstrapVersion.ID)

		// Query One
		queryBootstrapVersion := authRequest(`query ($id: BootstrapVersionID!) {
			getBootstrapVersion(bootstrapVersionID: $id) {
				id
				version
				satisfactory_version
				stability
				link
				changelog
				date
			}
		}`, token)
		queryBootstrapVersion.Var("id", createBootstrapVersionResponse.CreateBootstrapVersion.ID)

		var queryBootstrapVersionResponse struct {
			GetBootstrapVersion generated.BootstrapVersion
		}
		testza.AssertNoError(t, client.Run(ctx, queryBootstrapVersion, &queryBootstrapVersionResponse))
		testza.AssertEqual(t, createBootstrapVersionResponse.CreateBootstrapVersion.ID, queryBootstrapVersionResponse.GetBootstrapVersion.ID)
		testza.AssertEqual(t, version, queryBootstrapVersionResponse.GetBootstrapVersion.Version)
		testza.AssertEqual(t, 12345, queryBootstrapVersionResponse.GetBootstrapVersion.SatisfactoryVersion)
		testza.AssertEqual(t, generated.VersionStabilitiesBeta, queryBootstrapVersionResponse.GetBootstrapVersion.Stability)
		testza.AssertEqual(t, "example.com", queryBootstrapVersionResponse.GetBootstrapVersion.Link)
		testza.AssertEqual(t, "Hello World", queryBootstrapVersionResponse.GetBootstrapVersion.Changelog)

		// Update
		updateBootstrapVersion := authRequest(`mutation ($id: BootstrapVersionID!) {
			updateBootstrapVersion(
				bootstrapVersionId: $id,
				bootstrapVersion: {
					changelog: "Foo Bar",
				}
			) {
				id
			}
		}`, token)
		updateBootstrapVersion.Var("id", createBootstrapVersionResponse.CreateBootstrapVersion.ID)

		var updateBootstrapVersionResponse struct {
			UpdateBootstrapVersion generated.BootstrapVersion
		}
		testza.AssertNoError(t, client.Run(ctx, updateBootstrapVersion, &updateBootstrapVersionResponse))

		// Query Many
		queryBootstrapVersions := authRequest(`query {
			getBootstrapVersions {
				count
				bootstrap_versions {
					id
					version
					satisfactory_version
					stability
					link
					changelog
					date
				}
			}
		}`, token)

		var queryBootstrapVersionsResponse struct {
			GetBootstrapVersions generated.GetBootstrapVersions
		}
		testza.AssertNoError(t, client.Run(ctx, queryBootstrapVersions, &queryBootstrapVersionsResponse))
		testza.AssertEqual(t, 1, queryBootstrapVersionsResponse.GetBootstrapVersions.Count)
		testza.AssertEqual(t, 1, len(queryBootstrapVersionsResponse.GetBootstrapVersions.BootstrapVersions))
		testza.AssertEqual(t, createBootstrapVersionResponse.CreateBootstrapVersion.ID, queryBootstrapVersionsResponse.GetBootstrapVersions.BootstrapVersions[0].ID)
		testza.AssertEqual(t, version, queryBootstrapVersionsResponse.GetBootstrapVersions.BootstrapVersions[0].Version)
		testza.AssertEqual(t, 12345, queryBootstrapVersionsResponse.GetBootstrapVersions.BootstrapVersions[0].SatisfactoryVersion)
		testza.AssertEqual(t, generated.VersionStabilitiesBeta, queryBootstrapVersionsResponse.GetBootstrapVersions.BootstrapVersions[0].Stability)
		testza.AssertEqual(t, "example.com", queryBootstrapVersionsResponse.GetBootstrapVersions.BootstrapVersions[0].Link)
		testza.AssertEqual(t, "Foo Bar", queryBootstrapVersionsResponse.GetBootstrapVersions.BootstrapVersions[0].Changelog)

		// Delete
		deleteBootstrapVersion := authRequest(`mutation ($id: BootstrapVersionID!) {
			deleteBootstrapVersion(bootstrapVersionId: $id)
		}`, token)
		deleteBootstrapVersion.Var("id", createBootstrapVersionResponse.CreateBootstrapVersion.ID)

		var deleteBootstrapVersionResponse struct {
			DeleteBootstrapVersion bool
		}
		testza.AssertNoError(t, client.Run(ctx, deleteBootstrapVersion, &deleteBootstrapVersionResponse))
		testza.AssertTrue(t, deleteBootstrapVersionResponse.DeleteBootstrapVersion)
	}
}
