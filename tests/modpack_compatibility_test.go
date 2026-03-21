package tests

import (
	"testing"

	"github.com/MarvinJWendt/testza"

	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
)

func init() {
	config.SetConfigDir("../")
	db.EnableDebug()
}

func TestModpackCompatibility(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	tags := seedTags(ctx, t, token, client)
	mods := seedMods(ctx, t, token, client, tags[0])
	var objID string
	t.Run("Create", func(t *testing.T) {
		createRequest := authRequest(`mutation ($name: String!, $shortDescription: String!, $fullDescription: String!, $tags: [TagID!], $targets: [String!]!, $mods: [ModpackModInput!]!) {
					createModpack(modpack: {
						name: $name,
						short_description: $shortDescription,
						full_description: $fullDescription,
						tagIDs: $tags,
						targets: $targets,
						mods: $mods
					}) {
						id
						name
						short_description
						full_description
						creator_id
						hidden
						tags {
							id
							name
						}
						targets
						mods {
							mod_id
							version_constraint
						}
					}
				}`, token)
		createRequest.Var("name", "Test Modpack")
		createRequest.Var("shortDescription", "A test modpack for testing purposes")
		createRequest.Var("fullDescription", "This is a comprehensive test modpack that includes multiple mods for testing the modpack compatibility.")
		createRequest.Var("tags", tags)
		createRequest.Var("targets", []string{"Windows", "LinuxServer"})
		createRequest.Var("mods", []struct {
			ModID             string `json:"mod_id"`
			VersionConstraint string `json:"version_constraint"`
		}{
			{ModID: mods[0], VersionConstraint: ">=1.0.0"},
			{ModID: mods[1], VersionConstraint: "^2.0.0"},
		})

		var createResponse struct {
			CreateModpack generated.Modpack
		}
		err = client.Run(ctx, createRequest, &createResponse)
		testza.AssertNoError(t, err)
		objID = createResponse.CreateModpack.ID
	})

	t.Run("Working Compatibility", func(t *testing.T) {
		// First mod: EA Broken, EXP Works
		updateRequest1 := authRequest(`mutation ($id: ModID!) {
			updateMod(
				modId: $id,
				mod: {
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
			}}`, token)
		updateRequest1.Var("id", mods[0])
		err = client.Run(ctx, updateRequest1, nil)
		testza.AssertNoError(t, err)

		// Second mod: EA Works, EXP Damaged
		updateRequest2 := authRequest(`mutation ($id: ModID!) {
			updateMod(
				modId: $id,
				mod: {
					compatibility: {
						EA: {
							note: "Hello"
							state: Works
						}
						EXP: {
							note: "World",
							state: Damaged
						}
					}
				}
			) {
				id
			}}`, token)
		updateRequest2.Var("id", mods[1])
		err = client.Run(ctx, updateRequest2, nil)
		testza.AssertNoError(t, err)

		queryRequest := authRequest(`query ($modpackID: ModpackID!) {
				getModCompatibilities(modpackID: $modpackID) {
					worstEA {
						id
						name
						compatibility {
                			EA {
                    			state
              			  }
            			    EXP {
            			        state
            			    }
          				}
					}
					worstEXP {
						id
						name
						compatibility {
                			EA {
                    			state
              			  }
            			    EXP {
            			        state
            			    }
          				}
					}
				}
			}`, token)
		queryRequest.Var("modpackID", objID)

		var queryResponse struct {
			GetModCompatibilities generated.ModCompatibilities
		}
		err = client.Run(ctx, queryRequest, &queryResponse)
		testza.AssertNoError(t, err)

		testza.AssertEqual(t, 1, len(queryResponse.GetModCompatibilities.WorstEa))
		testza.AssertEqual(t, mods[0], queryResponse.GetModCompatibilities.WorstEa[0].ID)
		testza.AssertEqual(t, generated.CompatibilityStateBroken, queryResponse.GetModCompatibilities.WorstEa[0].Compatibility.Ea.State)

		testza.AssertEqual(t, 1, len(queryResponse.GetModCompatibilities.WorstExp))
		testza.AssertEqual(t, mods[1], queryResponse.GetModCompatibilities.WorstExp[0].ID)
		testza.AssertEqual(t, generated.CompatibilityStateDamaged, queryResponse.GetModCompatibilities.WorstExp[0].Compatibility.Exp.State)
	})
}
