package tests

import (
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

func TestGuides(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, userID, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	// Run Twice to detect any cache issues
	for i := 0; i < 2; i++ {
		// Create
		createGuide := authRequest(`mutation {
			createGuide(guide: {
				name: "Hello World",
				short_description: "Short description about the guide",
				guide: "The full guide text goes here."
			}) {
				id
			}
		}`, token)

		var createGuideResponse struct {
			CreateGuide generated.Guide
		}
		testza.AssertNoError(t, client.Run(ctx, createGuide, &createGuideResponse))
		testza.AssertNotEqual(t, "", createGuideResponse.CreateGuide.ID)

		// Query One
		queryGuide := authRequest(`query ($id: GuideID!) {
			getGuide(guideId: $id) {
				id
				name
				short_description
				guide
				user {
					id
				}
			}
		}`, token)
		queryGuide.Var("id", createGuideResponse.CreateGuide.ID)

		var queryGuideResponse struct {
			GetGuide generated.Guide
		}
		testza.AssertNoError(t, client.Run(ctx, queryGuide, &queryGuideResponse))
		testza.AssertEqual(t, createGuideResponse.CreateGuide.ID, queryGuideResponse.GetGuide.ID)
		testza.AssertEqual(t, "Hello World", queryGuideResponse.GetGuide.Name)
		testza.AssertEqual(t, "Short description about the guide", queryGuideResponse.GetGuide.ShortDescription)
		testza.AssertEqual(t, "The full guide text goes here.", queryGuideResponse.GetGuide.Guide)
		testza.AssertEqual(t, userID, queryGuideResponse.GetGuide.User.ID)

		// Update
		updateGuide := authRequest(`mutation ($id: GuideID!) {
			updateGuide(
				guideId: $id,
				guide: {
					name: "Foo Bar"
				}
			) {
				id
			}
		}`, token)
		updateGuide.Var("id", createGuideResponse.CreateGuide.ID)

		var updateGuideResponse struct {
			UpdateGuide generated.Guide
		}
		testza.AssertNoError(t, client.Run(ctx, updateGuide, &updateGuideResponse))

		// Query Many
		queryGuides := authRequest(`query {
			getGuides {
				count
				guides {
					id
					name
					short_description
					guide
					user {
						id
					}				
				}
			}
		}`, token)

		var queryGuidesResponse struct {
			GetGuides generated.GetGuides
		}
		testza.AssertNoError(t, client.Run(ctx, queryGuides, &queryGuidesResponse))
		testza.AssertEqual(t, 1, queryGuidesResponse.GetGuides.Count)
		testza.AssertEqual(t, 1, len(queryGuidesResponse.GetGuides.Guides))
		testza.AssertEqual(t, createGuideResponse.CreateGuide.ID, queryGuidesResponse.GetGuides.Guides[0].ID)
		testza.AssertEqual(t, "Foo Bar", queryGuidesResponse.GetGuides.Guides[0].Name)
		testza.AssertEqual(t, "Short description about the guide", queryGuidesResponse.GetGuides.Guides[0].ShortDescription)
		testza.AssertEqual(t, "The full guide text goes here.", queryGuidesResponse.GetGuides.Guides[0].Guide)
		testza.AssertEqual(t, userID, queryGuidesResponse.GetGuides.Guides[0].User.ID)

		// Delete
		deleteGuide := authRequest(`mutation ($id: GuideID!) {
			deleteGuide(guideId: $id)
		}`, token)
		deleteGuide.Var("id", createGuideResponse.CreateGuide.ID)

		var deleteGuideResponse struct {
			DeleteGuide bool
		}
		testza.AssertNoError(t, client.Run(ctx, deleteGuide, &deleteGuideResponse))
		testza.AssertTrue(t, deleteGuideResponse.DeleteGuide)
	}
}
