package tests

import (
	"strconv"
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

func TestModpackDB(t *testing.T) {
	ctx, _, stop := setup()
	defer stop()

	_, userID, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	testza.AssertNoError(t, db.From(ctx).Mod.Create().
		SetName("Fluffy Unicorns").
		SetModReference("FluffyUnicorns").
		SetShortDescription("A").
		SetFullDescription("B").
		SetLogo("C").
		SetCreatorID(userID).
		Exec(ctx))

	testza.AssertNoError(t, db.From(ctx).Mod.Create().
		SetName("Death Skulls").
		SetModReference("DeathSkulls").
		SetShortDescription("A").
		SetFullDescription("B").
		SetLogo("C").
		SetCreatorID(userID).
		Exec(ctx))

	testza.AssertNoError(t, db.From(ctx).Modpack.Create().
		SetName("Mega Pack").
		SetShortDescription("A").
		SetFullDescription("B").
		SetLogo("C").
		SetCreatorID(userID).
		Exec(ctx))
}

func TestModpacks(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, userID, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	// Second user for testing permissions
	_, _, err = makeUser(ctx)
	testza.AssertNoError(t, err)

	tags := seedTags(ctx, t, token, client)
	mods := seedMods(ctx, t, token, client, tags[0])

	// Run twice to detect any cache issues
	for i := range 2 {
		t.Run("Loop"+strconv.Itoa(i), func(t *testing.T) {
			var objID string
			var releaseID string

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
				createRequest.Var("name", "Test Modpack "+strconv.Itoa(i))
				createRequest.Var("shortDescription", "A test modpack for testing purposes")
				createRequest.Var("fullDescription", "This is a comprehensive test modpack that includes multiple mods for testing the modpack functionality.")
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
				testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
				testza.AssertNotEqual(t, "", createResponse.CreateModpack.ID)
				testza.AssertEqual(t, "Test Modpack "+strconv.Itoa(i), createResponse.CreateModpack.Name)
				testza.AssertEqual(t, "A test modpack for testing purposes", createResponse.CreateModpack.ShortDescription)
				testza.AssertEqual(t, userID, createResponse.CreateModpack.CreatorID)
				testza.AssertFalse(t, createResponse.CreateModpack.Hidden)
				testza.AssertEqual(t, 2, len(createResponse.CreateModpack.Tags))
				testza.AssertEqual(t, 2, len(createResponse.CreateModpack.Targets))
				testza.AssertEqual(t, 2, len(createResponse.CreateModpack.Mods))

				objID = createResponse.CreateModpack.ID
			})

			t.Run("Query One", func(t *testing.T) {
				queryRequest := authRequest(`query ($id: ModpackID!) {
					getModpack(modpackID: $id) {
						id
						name
						short_description
						full_description
						logo
						creator_id
						views
						installs
						hotness
						popularity
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
						releases {
							id
							version
						}
						parent {
							id
							name
						}
						children {
							id
							name
						}
					}
				}`, token)
				queryRequest.Var("id", objID)

				var queryResponse struct {
					GetModpack generated.Modpack
				}
				testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
				testza.AssertEqual(t, objID, queryResponse.GetModpack.ID)
				testza.AssertEqual(t, "Test Modpack "+strconv.Itoa(i), queryResponse.GetModpack.Name)
				testza.AssertEqual(t, "A test modpack for testing purposes", queryResponse.GetModpack.ShortDescription)
				testza.AssertEqual(t, userID, queryResponse.GetModpack.CreatorID)
				testza.AssertFalse(t, queryResponse.GetModpack.Hidden)
			})

			t.Run("Update", func(t *testing.T) {
				updateRequest := authRequest(`mutation ($id: ModpackID!, $name: String, $shortDescription: String, $fullDescription: String, $hidden: Boolean, $tags: [TagID!]) {
					updateModpack(modpackID: $id, modpack: {
						name: $name,
						short_description: $shortDescription,
						full_description: $fullDescription,
						hidden: $hidden,
						tagIDs: $tags
					}) {
						id
						name
						short_description
						full_description
						hidden
						tags {
							id
							name
						}
					}
				}`, token)
				updateRequest.Var("id", objID)
				updateRequest.Var("name", "Updated Modpack "+strconv.Itoa(i))
				updateRequest.Var("shortDescription", "Updated description")
				updateRequest.Var("fullDescription", "Updated full description with more details.")
				updateRequest.Var("hidden", false)
				updateRequest.Var("tags", []string{tags[0]}) // Only first tag

				var updateResponse struct {
					UpdateModpack generated.Modpack
				}
				testza.AssertNoError(t, client.Run(ctx, updateRequest, &updateResponse))
				testza.AssertEqual(t, objID, updateResponse.UpdateModpack.ID)
				testza.AssertEqual(t, "Updated Modpack "+strconv.Itoa(i), updateResponse.UpdateModpack.Name)
				testza.AssertEqual(t, "Updated description", updateResponse.UpdateModpack.ShortDescription)
				testza.AssertFalse(t, updateResponse.UpdateModpack.Hidden)
				testza.AssertEqual(t, 1, len(updateResponse.UpdateModpack.Tags))
			})

			t.Run("Create Release", func(t *testing.T) {
				createReleaseRequest := authRequest(`mutation ($modpackID: ModpackID!, $version: String!, $changelog: String!) {
					createModpackRelease(modpackID: $modpackID, release: {
						version: $version,
						changelog: $changelog
					}) {
						id
						version
						changelog
						lockfile
					}
				}`, token)
				createReleaseRequest.Var("modpackID", objID)
				createReleaseRequest.Var("version", "1.0."+strconv.Itoa(i))
				createReleaseRequest.Var("changelog", "Initial release for testing")

				var createReleaseResponse struct {
					CreateModpackRelease generated.ModpackRelease
				}
				testza.AssertNoError(t, client.Run(ctx, createReleaseRequest, &createReleaseResponse))
				testza.AssertNotEqual(t, "", createReleaseResponse.CreateModpackRelease.ID)
				testza.AssertEqual(t, "1.0."+strconv.Itoa(i), createReleaseResponse.CreateModpackRelease.Version)
				testza.AssertEqual(t, "Initial release for testing", createReleaseResponse.CreateModpackRelease.Changelog)
				testza.AssertNotEqual(t, "", createReleaseResponse.CreateModpackRelease.Lockfile)

				releaseID = createReleaseResponse.CreateModpackRelease.ID
			})

			t.Run("Query Release", func(t *testing.T) {
				queryReleaseRequest := authRequest(`query ($modpackID: ModpackID!, $version: String!) {
					getModpackRelease(modpackID: $modpackID, version: $version) {
						id
						version
						changelog
						lockfile
					}
				}`, token)
				queryReleaseRequest.Var("modpackID", objID)
				queryReleaseRequest.Var("version", "1.0."+strconv.Itoa(i))

				var queryReleaseResponse struct {
					GetModpackRelease generated.ModpackRelease
				}
				testza.AssertNoError(t, client.Run(ctx, queryReleaseRequest, &queryReleaseResponse))
				testza.AssertEqual(t, releaseID, queryReleaseResponse.GetModpackRelease.ID)
				testza.AssertEqual(t, "1.0."+strconv.Itoa(i), queryReleaseResponse.GetModpackRelease.Version)
				testza.AssertEqual(t, "Initial release for testing", queryReleaseResponse.GetModpackRelease.Changelog)
			})

			t.Run("Resolve Modpack", func(t *testing.T) {
				resolveRequest := authRequest(`mutation ($modpackID: ModpackID!, $targets: [String!]!) {
					resolveModpack(modpackID: $modpackID, targets: $targets)
				}`, token)
				resolveRequest.Var("modpackID", objID)
				resolveRequest.Var("targets", []string{"Windows"})

				var resolveResponse struct {
					ResolveModpack *string
				}
				testza.AssertNoError(t, client.Run(ctx, resolveRequest, &resolveResponse))
				testza.AssertNotNil(t, resolveResponse.ResolveModpack)
				testza.AssertNotEqual(t, "", *resolveResponse.ResolveModpack)
			})

			if i == 0 {
				t.Run("Delete Release", func(t *testing.T) {
					deleteReleaseRequest := authRequest(`mutation ($modpackID: ModpackID!, $version: String!) {
						deleteModpackRelease(modpackID: $modpackID, version: $version)
					}`, token)
					deleteReleaseRequest.Var("modpackID", objID)
					deleteReleaseRequest.Var("version", "1.0.0")

					var deleteReleaseResponse struct {
						DeleteModpackRelease bool
					}
					testza.AssertNoError(t, client.Run(ctx, deleteReleaseRequest, &deleteReleaseResponse))
					testza.AssertTrue(t, deleteReleaseResponse.DeleteModpackRelease)
				})

				t.Run("Delete", func(t *testing.T) {
					deleteRequest := authRequest(`mutation ($id: ModpackID!) {
						deleteModpack(modpackID: $id)
					}`, token)
					deleteRequest.Var("id", objID)

					var deleteResponse struct {
						DeleteModpack bool
					}
					testza.AssertNoError(t, client.Run(ctx, deleteRequest, &deleteResponse))
					testza.AssertTrue(t, deleteResponse.DeleteModpack)
				})
			}
		})
	}

	t.Run("Query Many", func(t *testing.T) {
		queryRequest := authRequest(`query {
			getModpacks(filter: {order: asc, order_by: created_at}) {
				count
				modpacks {
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
				}
			}
		}`, token)

		var queryResponse struct {
			GetModpacks generated.GetModpacks
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
		testza.AssertEqual(t, 1, queryResponse.GetModpacks.Count) // Only one left after deletion
		testza.AssertEqual(t, 1, len(queryResponse.GetModpacks.Modpacks))
		testza.AssertEqual(t, "Updated Modpack 1", queryResponse.GetModpacks.Modpacks[0].Name)
	})
}

func TestModpackRemix(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	tags := seedTags(ctx, t, token, client)
	mods := seedMods(ctx, t, token, client, tags[0])

	// Create parent modpack
	createParentRequest := authRequest(`mutation ($name: String!, $shortDescription: String!, $fullDescription: String!, $tags: [TagID!], $targets: [String!]!, $mods: [ModpackModInput!]!) {
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
		}
	}`, token)
	createParentRequest.Var("name", "Parent Modpack")
	createParentRequest.Var("shortDescription", "Original modpack")
	createParentRequest.Var("fullDescription", "This is the original modpack.")
	createParentRequest.Var("tags", tags)
	createParentRequest.Var("targets", []string{"Windows"})
	createParentRequest.Var("mods", []struct {
		ModID             string `json:"mod_id"`
		VersionConstraint string `json:"version_constraint"`
	}{
		{ModID: mods[0], VersionConstraint: ">=1.0.0"},
	})

	var createParentResponse struct {
		CreateModpack generated.Modpack
	}
	testza.AssertNoError(t, client.Run(ctx, createParentRequest, &createParentResponse))
	parentID := createParentResponse.CreateModpack.ID

	// Create child modpack (remix)
	createChildRequest := authRequest(`mutation ($name: String!, $shortDescription: String!, $fullDescription: String!, $parentID: ModpackID!, $tags: [TagID!], $targets: [String!]!, $mods: [ModpackModInput!]!) {
		createModpack(modpack: {
			name: $name,
			short_description: $shortDescription,
			full_description: $fullDescription,
			parent_id: $parentID,
			tagIDs: $tags,
			targets: $targets,
			mods: $mods
		}) {
			id
			name
			parent {
				id
				name
			}
		}
	}`, token)
	createChildRequest.Var("name", "Child Modpack")
	createChildRequest.Var("shortDescription", "Remix of original")
	createChildRequest.Var("fullDescription", "This is a remix of the original modpack.")
	createChildRequest.Var("parentID", parentID)
	createChildRequest.Var("tags", tags)
	createChildRequest.Var("targets", []string{"Windows"})
	createChildRequest.Var("mods", []struct {
		ModID             string `json:"mod_id"`
		VersionConstraint string `json:"version_constraint"`
	}{
		{ModID: mods[1], VersionConstraint: "^2.0.0"},
	})

	var createChildResponse struct {
		CreateModpack generated.Modpack
	}
	testza.AssertNoError(t, client.Run(ctx, createChildRequest, &createChildResponse))

	testza.AssertNotEqual(t, "", createChildResponse.CreateModpack.ID)
	testza.AssertEqual(t, "Child Modpack", createChildResponse.CreateModpack.Name)
	testza.AssertNotNil(t, createChildResponse.CreateModpack.Parent)
	testza.AssertEqual(t, parentID, createChildResponse.CreateModpack.Parent.ID)
	testza.AssertEqual(t, "Parent Modpack", createChildResponse.CreateModpack.Parent.Name)

	// Try to create a remix of a remix (should fail)
	createGrandchildRequest := authRequest(`mutation ($name: String!, $shortDescription: String!, $fullDescription: String!, $parentID: ModpackID!, $tags: [TagID!], $targets: [String!]!, $mods: [ModpackModInput!]!) {
		createModpack(modpack: {
			name: $name,
			short_description: $shortDescription,
			full_description: $fullDescription,
			parent_id: $parentID,
			tagIDs: $tags,
			targets: $targets,
			mods: $mods
		}) {
			id
		}
	}`, token)
	createGrandchildRequest.Var("name", "Grandchild Modpack")
	createGrandchildRequest.Var("shortDescription", "Remix of remix")
	createGrandchildRequest.Var("fullDescription", "This should fail.")
	createGrandchildRequest.Var("parentID", createChildResponse.CreateModpack.ID)
	createGrandchildRequest.Var("tags", tags)
	createGrandchildRequest.Var("targets", []string{"Windows"})
	createGrandchildRequest.Var("mods", []struct {
		ModID             string `json:"mod_id"`
		VersionConstraint string `json:"version_constraint"`
	}{
		{ModID: mods[0], VersionConstraint: ">=1.0.0"},
	})

	var createGrandchildResponse struct {
		CreateModpack generated.Modpack
	}
	err = client.Run(ctx, createGrandchildRequest, &createGrandchildResponse)
	testza.AssertNotNil(t, err)
	testza.AssertContains(t, err.Error(), "cannot create a remix of a remix")

	// Query parent to see children
	queryParentRequest := authRequest(`query ($id: ModpackID!) {
		getModpack(modpackID: $id) {
			id
			name
			children {
				id
				name
			}
		}
	}`, token)
	queryParentRequest.Var("id", parentID)

	var queryParentResponse struct {
		GetModpack generated.Modpack
	}
	testza.AssertNoError(t, client.Run(ctx, queryParentRequest, &queryParentResponse))
	testza.AssertEqual(t, parentID, queryParentResponse.GetModpack.ID)
	testza.AssertEqual(t, 1, len(queryParentResponse.GetModpack.Children))
	testza.AssertEqual(t, "Child Modpack", queryParentResponse.GetModpack.Children[0].Name)
}

func TestMyModpack(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, userID, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	token2, userID2, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	token3, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	tags := seedTags(ctx, t, token, client)
	mods := seedMods(ctx, t, token, client, tags[0])

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
		createRequest.Var("name", "Test Modpack 1")
		createRequest.Var("shortDescription", "A test modpack 1 for testing purposes")
		createRequest.Var("fullDescription", "This is a comprehensive test modpack 1 that includes multiple mods for testing the modpack functionality.")
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
		testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
		testza.AssertNotEqual(t, "", createResponse.CreateModpack.ID)
		testza.AssertEqual(t, "Test Modpack 1", createResponse.CreateModpack.Name)
		testza.AssertEqual(t, "A test modpack 1 for testing purposes", createResponse.CreateModpack.ShortDescription)
		testza.AssertEqual(t, userID, createResponse.CreateModpack.CreatorID)
		testza.AssertFalse(t, createResponse.CreateModpack.Hidden)
		testza.AssertEqual(t, 2, len(createResponse.CreateModpack.Tags))
		testza.AssertEqual(t, 2, len(createResponse.CreateModpack.Targets))
		testza.AssertEqual(t, 2, len(createResponse.CreateModpack.Mods))

		createRequest2 := authRequest(`mutation ($name: String!, $shortDescription: String!, $fullDescription: String!, $tags: [TagID!], $targets: [String!]!, $mods: [ModpackModInput!]!) {
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
		}`, token2)
		createRequest2.Var("name", "Test Modpack 2")
		createRequest2.Var("shortDescription", "A test modpack 2 for testing purposes")
		createRequest2.Var("fullDescription", "This is a comprehensive test modpack 2 that includes multiple mods for testing the modpack functionality.")
		createRequest2.Var("tags", tags)
		createRequest2.Var("targets", []string{"Windows", "LinuxServer"})
		createRequest2.Var("mods", []struct {
			ModID             string `json:"mod_id"`
			VersionConstraint string `json:"version_constraint"`
		}{
			{ModID: mods[2], VersionConstraint: ">=1.0.0"},
			{ModID: mods[3], VersionConstraint: "^2.0.0"},
		})

		var createResponse2 struct {
			CreateModpack generated.Modpack
		}
		testza.AssertNoError(t, client.Run(ctx, createRequest2, &createResponse2))
		testza.AssertNotEqual(t, "", createResponse2.CreateModpack.ID)
		testza.AssertEqual(t, "Test Modpack 2", createResponse2.CreateModpack.Name)
		testza.AssertEqual(t, "A test modpack 2 for testing purposes", createResponse2.CreateModpack.ShortDescription)
		testza.AssertEqual(t, userID2, createResponse2.CreateModpack.CreatorID)
		testza.AssertFalse(t, createResponse2.CreateModpack.Hidden)
		testza.AssertEqual(t, 2, len(createResponse2.CreateModpack.Tags))
		testza.AssertEqual(t, 2, len(createResponse2.CreateModpack.Targets))
		testza.AssertEqual(t, 2, len(createResponse2.CreateModpack.Mods))
	})

	t.Run("Query Many", func(t *testing.T) {
		queryRequest := authRequest(`query {
			getModpacks(filter: {order: asc, order_by: created_at}) {
				count
				modpacks {
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
				}
			}
		}`, token)

		var queryResponse struct {
			GetModpacks generated.GetModpacks
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
		testza.AssertEqual(t, 2, queryResponse.GetModpacks.Count)
		testza.AssertEqual(t, 2, len(queryResponse.GetModpacks.Modpacks))
		testza.AssertEqual(t, "Test Modpack 1", queryResponse.GetModpacks.Modpacks[0].Name)
	})

	t.Run("Query Many", func(t *testing.T) {
		queryRequest := authRequest(`query {
			getModpacks(filter: {order: asc, order_by: created_at}) {
				count
				modpacks {
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
				}
			}
		}`, token2)

		var queryResponse struct {
			GetModpacks generated.GetModpacks
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
		testza.AssertEqual(t, 2, queryResponse.GetModpacks.Count)
		testza.AssertEqual(t, 2, len(queryResponse.GetModpacks.Modpacks))
		testza.AssertEqual(t, "Test Modpack 1", queryResponse.GetModpacks.Modpacks[0].Name)
	})

	t.Run("Query Own 1", func(t *testing.T) {
		queryRequest := authRequest(`query {
			getMyModpacks(filter: {order: asc, order_by: created_at}) {
				count
				modpacks {
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
				}
			}
		}`, token)

		var queryResponse struct {
			GetMyModpacks generated.GetMyModpacks
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
		testza.AssertEqual(t, 1, queryResponse.GetMyModpacks.Count)
		testza.AssertEqual(t, 1, len(queryResponse.GetMyModpacks.Modpacks))
		testza.AssertEqual(t, "Test Modpack 1", queryResponse.GetMyModpacks.Modpacks[0].Name)
	})

	t.Run("Query Own 2", func(t *testing.T) {
		queryRequest := authRequest(`query {
			getMyModpacks(filter: {order: asc, order_by: created_at}) {
				count
				modpacks {
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
				}
			}
		}`, token2)

		var queryResponse struct {
			GetMyModpacks generated.GetMyModpacks
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
		testza.AssertEqual(t, 1, queryResponse.GetMyModpacks.Count)
		testza.AssertEqual(t, 1, len(queryResponse.GetMyModpacks.Modpacks))
		testza.AssertEqual(t, "Test Modpack 2", queryResponse.GetMyModpacks.Modpacks[0].Name)
	})

	t.Run("Query Own 3", func(t *testing.T) {
		queryRequest := authRequest(`query {
			getMyModpacks(filter: {order: asc, order_by: created_at}) {
				count
				modpacks {
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
				}
			}
		}`, token3)

		var queryResponse struct {
			GetMyModpacks generated.GetMyModpacks
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
		testza.AssertEqual(t, 0, queryResponse.GetMyModpacks.Count)
		testza.AssertEqual(t, 0, len(queryResponse.GetMyModpacks.Modpacks))
	})

	t.Run("GetUser Modpacks", func(t *testing.T) {
		req := authRequest(`query ($user: UserID!) {
		getUser(userId: $user) {
			id
			modpacks {
				modpack {
					id
					name
				}
			}
		}
	}`, token)

		req.Var("user", userID)

		var resp struct {
			GetUser struct {
				Modpacks []struct {
					Modpack struct {
						ID   string
						Name string
					}
				}
			}
		}

		testza.AssertNoError(t, client.Run(ctx, req, &resp))

		testza.AssertEqual(t, 1, len(resp.GetUser.Modpacks))
		testza.AssertEqual(t, "Test Modpack 1", resp.GetUser.Modpacks[0].Modpack.Name)
	})

	t.Run("GetMyModpacks Unauthorized", func(t *testing.T) {
		req := authRequest(`query {
		getMyModpacks {
			count
		}
	}`, "") // no token

		var resp struct{}

		err := client.Run(ctx, req, &resp)
		testza.AssertNotNil(t, err)
		testza.AssertContains(t, err.Error(), "graphql: user not logged in")
	})
}
