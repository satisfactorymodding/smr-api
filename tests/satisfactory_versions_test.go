package tests

import (
	"strconv"
	"testing"

	"github.com/MarvinJWendt/testza"

	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/migrations"
)

func init() {
	migrations.SetMigrationDir("../migrations")
	config.SetConfigDir("../")
	db.EnableDebug()
}

func TestSatisfactoryVersions(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	secondaryVersions := [2]int{
		234567,
		345678,
	}

	// Run Twice to detect any cache issues
	for i := range 2 {
		t.Run("Loop"+strconv.Itoa(i), func(t *testing.T) {
			var objID string

			t.Run("Create", func(t *testing.T) {
				createRequest := authRequest(`mutation {
				  createSatisfactoryVersion(input: {
					version: 123456,
					engine_version: "5.1"
				  }) {
					id
				  }
				}`, token)

				var createResponse struct {
					CreateSatisfactoryVersion generated.SatisfactoryVersion
				}
				testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))

				testza.AssertNotEqual(t, "", createResponse.CreateSatisfactoryVersion.ID)

				objID = createResponse.CreateSatisfactoryVersion.ID
			})

			t.Run("Query One", func(t *testing.T) {
				queryRequest := authRequest(`query ($id: SatisfactoryVersionID!) {
				  getSatisfactoryVersion(id: $id) {
					id
					version
					engine_version
				  }
				}`, token)
				queryRequest.Var("id", objID)

				var queryResponse struct {
					GetSatisfactoryVersion generated.SatisfactoryVersion
				}
				testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))

				testza.AssertEqual(t, objID, queryResponse.GetSatisfactoryVersion.ID)
				testza.AssertEqual(t, 123456, queryResponse.GetSatisfactoryVersion.Version)
				testza.AssertEqual(t, "5.1", queryResponse.GetSatisfactoryVersion.EngineVersion)
			})

			t.Run("Update", func(t *testing.T) {
				updateRequest := authRequest(`mutation ($id: SatisfactoryVersionID!, $version: Int!) {
				  updateSatisfactoryVersion(
					id: $id,
					input: {
					  version: $version
					}
				  ) {
					id
				  }
				}`, token)
				updateRequest.Var("id", objID)
				updateRequest.Var("version", secondaryVersions[i])

				var updateResponse struct {
					UpdateSatisfactoryVersion generated.SatisfactoryVersion
				}
				testza.AssertNoError(t, client.Run(ctx, updateRequest, &updateResponse))
			})

			t.Run("Query Many", func(t *testing.T) {
				queryRequest := authRequest(`{
				  getSatisfactoryVersions {
				    id
				    version
				    engine_version
				  }
				}`, token)

				var queryResponse struct {
					GetSatisfactoryVersions []generated.SatisfactoryVersion
				}
				testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))

				testza.AssertEqual(t, 1, len(queryResponse.GetSatisfactoryVersions))
				testza.AssertEqual(t, objID, queryResponse.GetSatisfactoryVersions[0].ID)
				testza.AssertEqual(t, secondaryVersions[i], queryResponse.GetSatisfactoryVersions[0].Version)
				testza.AssertEqual(t, "5.1", queryResponse.GetSatisfactoryVersions[0].EngineVersion)
			})

			t.Run("Delete", func(t *testing.T) {
				deleteRequest := authRequest(`mutation ($id: SatisfactoryVersionID!) {
				  deleteSatisfactoryVersion(id: $id)
				}`, token)
				deleteRequest.Var("id", objID)

				var deleteResponse struct {
					DeleteSatisfactoryVersion bool
				}
				testza.AssertNoError(t, client.Run(ctx, deleteRequest, &deleteResponse))

				testza.AssertTrue(t, deleteResponse.DeleteSatisfactoryVersion)
			})
		})
	}
}
