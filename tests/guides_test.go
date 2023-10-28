package tests

import (
	"strconv"
	"testing"

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

func TestGuides(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, userID, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	// Run Twice to detect any cache issues
	for i := 0; i < 2; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var guideID string

			t.Run("Create", func(t *testing.T) {
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

				guideID = createGuideResponse.CreateGuide.ID
			})

			t.Run("Query One", func(t *testing.T) {
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
				queryGuide.Var("id", guideID)

				var queryGuideResponse struct {
					GetGuide generated.Guide
				}
				testza.AssertNoError(t, client.Run(ctx, queryGuide, &queryGuideResponse))
				testza.AssertEqual(t, guideID, queryGuideResponse.GetGuide.ID)
				testza.AssertEqual(t, "Hello World", queryGuideResponse.GetGuide.Name)
				testza.AssertEqual(t, "Short description about the guide", queryGuideResponse.GetGuide.ShortDescription)
				testza.AssertEqual(t, "The full guide text goes here.", queryGuideResponse.GetGuide.Guide)
				testza.AssertEqual(t, userID, queryGuideResponse.GetGuide.User.ID)
			})

			t.Run("Update", func(t *testing.T) {
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
				updateGuide.Var("id", guideID)

				var updateGuideResponse struct {
					UpdateGuide generated.Guide
				}
				testza.AssertNoError(t, client.Run(ctx, updateGuide, &updateGuideResponse))
			})

			t.Run("Query Many", func(t *testing.T) {
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
				testza.AssertEqual(t, guideID, queryGuidesResponse.GetGuides.Guides[0].ID)
				testza.AssertEqual(t, "Foo Bar", queryGuidesResponse.GetGuides.Guides[0].Name)
				testza.AssertEqual(t, "Short description about the guide", queryGuidesResponse.GetGuides.Guides[0].ShortDescription)
				testza.AssertEqual(t, "The full guide text goes here.", queryGuidesResponse.GetGuides.Guides[0].Guide)
				testza.AssertEqual(t, userID, queryGuidesResponse.GetGuides.Guides[0].User.ID)
			})

			t.Run("Delete", func(t *testing.T) {
				deleteGuide := authRequest(`mutation ($id: GuideID!) {
					deleteGuide(guideId: $id)
				}`, token)
				deleteGuide.Var("id", guideID)

				var deleteGuideResponse struct {
					DeleteGuide bool
				}
				testza.AssertNoError(t, client.Run(ctx, deleteGuide, &deleteGuideResponse))
				testza.AssertTrue(t, deleteGuideResponse.DeleteGuide)
			})
		})
	}
}
