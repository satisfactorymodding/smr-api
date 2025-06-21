package tests

import (
	"strconv"
	"testing"

	"github.com/MarvinJWendt/testza"

	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/util"
)

func init() {
	config.SetConfigDir("../")
	db.EnableDebug()
}

func TestModRateLimit(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	tags := seedTags(ctx, t, token, client)

	// Create util.ModsPer24h mods to fill the rate limit
	for i := range util.ModsPer24h {
		createRequest := authRequest(`mutation ($mod_reference: ModReference!, $tags: [TagID!]) {
			createMod(mod: {
				name: "Rate Limit Test Mod",
				short_description: "Testing rate limiting functionality",
				full_description: "Lorem ipsum dolor sit amet",
				mod_reference: $mod_reference,
				tagIDs: $tags
			}) {
				id
				name
				mod_reference
			}
		}`, token)
		createRequest.Var("mod_reference", "ratelimit"+strconv.Itoa(i))
		createRequest.Var("tags", []string{tags[0]})

		var createResponse struct {
			CreateMod generated.Mod
		}
		err := client.Run(ctx, createRequest, &createResponse)
		testza.AssertNoError(t, err)
		testza.AssertNotNil(t, createResponse.CreateMod)
	}

	// Should fail with rate limit
	createRequest := authRequest(`mutation ($mod_reference: ModReference!, $tags: [TagID!]) {
		createMod(mod: {
			name: "Should Fail Mod",
			short_description: "This should fail due to rate limiting",
			full_description: "Lorem ipsum dolor sit amet",
			mod_reference: $mod_reference,
			tagIDs: $tags
		}) {
			id
			name
			mod_reference
		}
	}`, token)
	createRequest.Var("mod_reference", "shouldfail")
	createRequest.Var("tags", []string{tags[0]})

	var createResponse struct {
		CreateMod generated.Mod
	}
	err = client.Run(ctx, createRequest, &createResponse)
	testza.AssertNotNil(t, err)
	if err != nil {
		testza.AssertContains(t, err.Error(), "please wait")
		testza.AssertContains(t, err.Error(), "minutes to post another mod")
		// Verify that the error message does not contain negative minutes
		testza.AssertNotContains(t, err.Error(), "-")
	}
}

func TestMods(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, userID, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	_, userID2, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	tags := seedTags(ctx, t, token, client)

	// Run Twice to detect any cache issues
	for i := range 2 {
		t.Run("Loop"+strconv.Itoa(i), func(t *testing.T) {
			var objID string

			modReference := "hello" + strconv.Itoa(i)

			t.Run("Create", func(t *testing.T) {
				createRequest := authRequest(`mutation ($mod_reference: ModReference!, $tags: [TagID!]) {
					createMod(mod: {
						name: "Hello World",
						short_description: "Foo Bar 123 Foo Bar 123",
						full_description: "Lorem ipsum dolor sit amet",
						mod_reference: $mod_reference,
						tagIDs: $tags,
						toggle_network_use: true,
						toggle_explicit_content: true
					}) {
						id
					}
				}`, token)
				createRequest.Var("mod_reference", modReference)
				createRequest.Var("tags", tags)

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
						logo
						source_url
						creator_id
						approved
						views
						downloads
						hotness
						popularity
						updated_at
						created_at
						last_version_date
						mod_reference
						hidden
						toggle_network_use
						toggle_explicit_content
						tags {
						  id
						  name
						  description
						}
						compatibility {
						  EA {
							state
							note
						  }
						  EXP {
							state
							note
						  }
						}
						authors {
						  user_id
						  role
						}
						latestVersions {
						  alpha {
							id
						  }
						  beta {
							id
						  }
						  release {
							id
						  }
						}
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
				testza.AssertTrue(t, queryResponse.GetMod.ToggleNetworkUse)
				testza.AssertTrue(t, queryResponse.GetMod.ToggleExplicitContent)
			})

			t.Run("Query One By Reference", func(t *testing.T) {
				queryRequest := authRequest(`query ($modReference: ModReference!) {
					getModByReference(modReference: $modReference) {
						id
						name
						short_description
						full_description
						mod_reference
						creator_id
					}
				}`, token)
				queryRequest.Var("modReference", modReference)

				var queryResponse struct {
					GetModByReference generated.Mod
				}
				testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
				testza.AssertEqual(t, objID, queryResponse.GetModByReference.ID)
				testza.AssertEqual(t, "Hello World", queryResponse.GetModByReference.Name)
				testza.AssertEqual(t, "Foo Bar 123 Foo Bar 123", queryResponse.GetModByReference.ShortDescription)
				fullDescription := "Lorem ipsum dolor sit amet"
				testza.AssertEqual(t, &fullDescription, queryResponse.GetModByReference.FullDescription)
				testza.AssertEqual(t, modReference, queryResponse.GetModByReference.ModReference)
				testza.AssertEqual(t, userID, queryResponse.GetModByReference.CreatorID)
			})

			t.Run("Query One By ID or Reference", func(t *testing.T) {
				for _, s := range []string{objID, modReference} {
					queryRequest := authRequest(`query ($modIdOrReference: String!) {
						getModByIdOrReference(modIdOrReference: $modIdOrReference) {
							id
							name
							short_description
							full_description
							mod_reference
							creator_id
						}
					}`, token)
					queryRequest.Var("modIdOrReference", s)

					var queryResponse struct {
						GetModByIDOrReference generated.Mod
					}
					testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
					testza.AssertEqual(t, objID, queryResponse.GetModByIDOrReference.ID)
					testza.AssertEqual(t, "Hello World", queryResponse.GetModByIDOrReference.Name)
					testza.AssertEqual(t, "Foo Bar 123 Foo Bar 123", queryResponse.GetModByIDOrReference.ShortDescription)
					fullDescription := "Lorem ipsum dolor sit amet"
					testza.AssertEqual(t, &fullDescription, queryResponse.GetModByIDOrReference.FullDescription)
					testza.AssertEqual(t, modReference, queryResponse.GetModByIDOrReference.ModReference)
					testza.AssertEqual(t, userID, queryResponse.GetModByIDOrReference.CreatorID)
				}
			})

			t.Run("Update", func(t *testing.T) {
				updateRequest := authRequest(`mutation ($id: ModID!, $tags: [TagID!], $authors: [UpdateUserMod!]) {
					updateMod(
						modId: $id,
						mod: {
							name: "Foo Bar",
							tagIDs: $tags,
							authors: $authors,
							toggle_network_use: false,
							toggle_explicit_content: false,
							compatibility: {
								EA: {
									note: "Hello"
									state: Broken
								}
								EXP: {
									note: "World",
									state: Works
								}
							}
						}
					) {
						id
					}
				}`, token)
				updateRequest.Var("id", objID)
				updateRequest.Var("tags", tags)
				updateRequest.Var("authors", []struct {
					Role   string `json:"role"`
					UserID string `json:"user_id"`
				}{
					{
						Role:   "creator",
						UserID: userID,
					},
					{
						Role:   "editor",
						UserID: userID2,
					},
				})

				var updateResponse struct {
					UpdateMod generated.Mod
				}
				testza.AssertNoError(t, client.Run(ctx, updateRequest, &updateResponse))
			})

			t.Run("Query Many", func(t *testing.T) {
				queryRequest := authRequest(`query {
					getMods(filter: {order: asc, order_by: created_at}) {
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
				testza.AssertEqual(t, 2, queryResponse.GetMods.Count)
				testza.AssertEqual(t, 2, len(queryResponse.GetMods.Mods))
				testza.AssertEqual(t, objID, queryResponse.GetMods.Mods[1].ID)
				testza.AssertEqual(t, "Foo Bar", queryResponse.GetMods.Mods[1].Name)
				testza.AssertEqual(t, "Foo Bar 123 Foo Bar 123", queryResponse.GetMods.Mods[1].ShortDescription)
				fullDescription := "Lorem ipsum dolor sit amet"
				testza.AssertEqual(t, &fullDescription, queryResponse.GetMods.Mods[1].FullDescription)
				testza.AssertEqual(t, modReference, queryResponse.GetMods.Mods[1].ModReference)
				testza.AssertEqual(t, userID, queryResponse.GetMods.Mods[1].CreatorID)
				testza.AssertFalse(t, queryResponse.GetMods.Mods[1].ToggleNetworkUse)
				testza.AssertFalse(t, queryResponse.GetMods.Mods[1].ToggleExplicitContent)
			})

			t.Run("Query My Mods", func(t *testing.T) {
				queryRequest := authRequest(`query {
					getMyMods(filter: {order: asc, order_by: created_at}) {
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
					GetMyMods generated.GetMyMods
				}
				testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
				testza.AssertEqual(t, 2, queryResponse.GetMyMods.Count)
				testza.AssertEqual(t, 2, len(queryResponse.GetMyMods.Mods))
				testza.AssertEqual(t, objID, queryResponse.GetMyMods.Mods[1].ID)
				testza.AssertEqual(t, "Foo Bar", queryResponse.GetMyMods.Mods[1].Name)
				testza.AssertEqual(t, "Foo Bar 123 Foo Bar 123", queryResponse.GetMyMods.Mods[1].ShortDescription)
				fullDescription := "Lorem ipsum dolor sit amet"
				testza.AssertEqual(t, &fullDescription, queryResponse.GetMyMods.Mods[1].FullDescription)
				testza.AssertEqual(t, modReference, queryResponse.GetMyMods.Mods[1].ModReference)
				testza.AssertEqual(t, userID, queryResponse.GetMyMods.Mods[1].CreatorID)
			})

			if i == 0 {
				t.Run("Approve", func(t *testing.T) {
					approveRequest := authRequest(`mutation ApproveMod($id: ModID!) {
					  approveMod(modId: $id)
					}`, token)
					approveRequest.Var("id", objID)

					var approveResponse struct {
						ApproveMod bool
					}
					testza.AssertNoError(t, client.Run(ctx, approveRequest, &approveResponse))
					testza.AssertTrue(t, approveResponse.ApproveMod)
				})
			} else {
				t.Run("Deny", func(t *testing.T) {
					denyRequest := authRequest(`mutation DenyMod($id: ModID!) {
					  denyMod(modId: $id)
					}`, token)
					denyRequest.Var("id", objID)

					var denyResponse struct {
						DenyMod bool
					}
					testza.AssertNoError(t, client.Run(ctx, denyRequest, &denyResponse))
					testza.AssertTrue(t, denyResponse.DenyMod)
				})
			}

			if i == 0 {
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
			}
		})
	}

	t.Run("Prevent Mod Reference Matching Existing Mod ID", func(t *testing.T) {
		// First create a mod to get its ID
		token, _, err := makeUser(ctx)
		testza.AssertNoError(t, err)

		// Create the first mod
		createRequest := authRequest(`mutation ($mod_reference: ModReference!) {
			createMod(mod: {
				name: "First Mod",
				short_description: "This is the first mod for testing",
				full_description: "A full description for the first mod",
				mod_reference: $mod_reference
			}) {
				id
			}
		}`, token)
		createRequest.Var("mod_reference", "first-mod")

		var createResponse struct {
			CreateMod generated.Mod
		}
		testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
		testza.AssertNotEqual(t, "", createResponse.CreateMod.ID)

		existingModID := createResponse.CreateMod.ID

		// Now try to create another mod with a reference that matches the first mod's ID
		// This should fail
		createRequest2 := authRequest(`mutation ($mod_reference: ModReference!) {
			createMod(mod: {
				name: "Second Mod",
				short_description: "This should fail because reference matches existing mod ID",
				full_description: "A full description for the second mod",
				mod_reference: $mod_reference
			}) {
				id
			}
		}`, token)
		createRequest2.Var("mod_reference", existingModID) // Using existing mod ID as reference

		var createResponse2 struct {
			CreateMod generated.Mod
		}

		err = client.Run(ctx, createRequest2, &createResponse2)
		testza.AssertNotNil(t, err)
		testza.AssertContains(t, err.Error(), "mod reference cannot match an existing mod ID")
	})
}
