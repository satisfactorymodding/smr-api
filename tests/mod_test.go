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

// TODO Add mod tag test
// TODO Add rate limit test

func TestMods(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, userID, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	// Run Twice to detect any cache issues
	for i := 0; i < 2; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var objID string

			modReference := "hello" + strconv.Itoa(i)

			t.Run("Create", func(t *testing.T) {
				createRequest := authRequest(`mutation ($mod_reference: ModReference!) {
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

				objID = createResponse.CreateMod.ID
			})

			t.Run("Query One", func(t *testing.T) {
				queryRequest := authRequest(`query ($id: ModID!) {
					getMod(modId: $id) {
						id
						name
						short_description
						full_description
						mod_reference
						creator_id
					}
				}`, token)
				queryRequest.Var("id", objID)

				var queryResponse struct {
					GetMod generated.Mod
				}
				testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
				testza.AssertEqual(t, objID, queryResponse.GetMod.ID)
				testza.AssertEqual(t, "Hello World", queryResponse.GetMod.Name)
				testza.AssertEqual(t, "Foo Bar 123 Foo Bar 123", queryResponse.GetMod.ShortDescription)
				fullDescription := "Lorem ipsum dolor sit amet"
				testza.AssertEqual(t, &fullDescription, queryResponse.GetMod.FullDescription)
				testza.AssertEqual(t, modReference, queryResponse.GetMod.ModReference)
				testza.AssertEqual(t, userID, queryResponse.GetMod.CreatorID)
			})

			t.Run("Update", func(t *testing.T) {
				updateRequest := authRequest(`mutation ($id: ModID!) {
					updateMod(
						modId: $id,
						mod: {
							name: "Foo Bar"
						}
					) {
						id
					}
				}`, token)
				updateRequest.Var("id", objID)

				var updateResponse struct {
					UpdateMod generated.Mod
				}
				testza.AssertNoError(t, client.Run(ctx, updateRequest, &updateResponse))
			})

			t.Run("Query Many", func(t *testing.T) {
				queryRequest := authRequest(`query {
					getMods {
						count
						mods {
							id
							name
							short_description
							full_description
							mod_reference
							creator_id		
						}
					}
				}`, token)

				var queryResponse struct {
					GetMods generated.GetMods
				}
				testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
				testza.AssertEqual(t, 1, queryResponse.GetMods.Count)
				testza.AssertEqual(t, 1, len(queryResponse.GetMods.Mods))
				testza.AssertEqual(t, objID, queryResponse.GetMods.Mods[0].ID)
				testza.AssertEqual(t, "Foo Bar", queryResponse.GetMods.Mods[0].Name)
				testza.AssertEqual(t, "Foo Bar 123 Foo Bar 123", queryResponse.GetMods.Mods[0].ShortDescription)
				fullDescription := "Lorem ipsum dolor sit amet"
				testza.AssertEqual(t, &fullDescription, queryResponse.GetMods.Mods[0].FullDescription)
				testza.AssertEqual(t, modReference, queryResponse.GetMods.Mods[0].ModReference)
				testza.AssertEqual(t, userID, queryResponse.GetMods.Mods[0].CreatorID)
			})

			t.Run("Delete", func(t *testing.T) {
				deleteRequest := authRequest(`mutation ($id: ModID!) {
					deleteMod(modId: $id)
				}`, token)
				deleteRequest.Var("id", objID)

				var deleteResponse struct {
					DeleteMod bool
				}
				testza.AssertNoError(t, client.Run(ctx, deleteRequest, &deleteResponse))
				testza.AssertTrue(t, deleteResponse.DeleteMod)
			})
		})
	}
}
