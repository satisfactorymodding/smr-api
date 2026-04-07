package tests

import (
	"context"
	"os"
	"testing"

	"github.com/MarvinJWendt/testza"
	"github.com/machinebox/graphql"

	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/ent"
)

func init() {
	os.Setenv("NO_COLOR", "1")
	config.SetConfigDir("../")
	db.EnableDebug()
}

// $env:NO_COLOR=1; go test -v ./tests -run TestGetModpackTargetSupport
func TestGetModpackTargetSupport(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, userID, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	tags := seedTags(ctx, t, token, client)
	modIDs := seedMods(ctx, t, token, client, tags[0])

	t.Run("Modpack with all targets supported", func(t *testing.T) {
		modpackID := createTestModpack(ctx, t, client, token, userID, []string{modIDs[0]})

		createVersionsForTargets(ctx, t, modIDs[0], true, []string{"Windows", "WindowsServer", "LinuxServer"})

		queryRequest := authRequest(`query ($modpackID: ModpackID!) {
			getModpackTargetSupport(modpackID: $modpackID) {
				target_name
			}
		}`, token)
		queryRequest.Var("modpackID", modpackID)

		var response struct {
			GetModpackTargetSupport []*struct {
				TargetName string `json:"target_name"`
			}
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &response))
		testza.AssertEqual(t, 3, len(response.GetModpackTargetSupport))

		targetNames := make(map[string]bool)
		for _, target := range response.GetModpackTargetSupport {
			targetNames[target.TargetName] = true
		}
		testza.AssertTrue(t, targetNames["Windows"])
		testza.AssertTrue(t, targetNames["WindowsServer"])
		testza.AssertTrue(t, targetNames["LinuxServer"])
	})

	t.Run("Modpack with missing target in one mod", func(t *testing.T) {
		modpackID := createTestModpack(ctx, t, client, token, userID, []string{modIDs[0], modIDs[1]})

		createVersionsForTargets(ctx, t, modIDs[0], true, []string{"Windows", "WindowsServer", "LinuxServer"})
		createVersionsForTargets(ctx, t, modIDs[1], true, []string{"Windows"})

		queryRequest := authRequest(`query ($modpackID: ModpackID!) {
			getModpackTargetSupport(modpackID: $modpackID) {
				target_name
			}
		}`, token)
		queryRequest.Var("modpackID", modpackID)

		var response struct {
			GetModpackTargetSupport []*struct {
				TargetName string `json:"target_name"`
			}
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &response))
		testza.AssertEqual(t, 1, len(response.GetModpackTargetSupport))

		targetNames := make(map[string]bool)
		for _, target := range response.GetModpackTargetSupport {
			targetNames[target.TargetName] = true
		}
		testza.AssertTrue(t, targetNames["Windows"])
		testza.AssertFalse(t, targetNames["WindowsServer"])
		testza.AssertFalse(t, targetNames["LinuxServer"])
	})

	t.Run("Modpack with non-required missing targets", func(t *testing.T) {
		modpackID := createTestModpack(ctx, t, client, token, userID, []string{modIDs[0], modIDs[2]})

		createVersionsForTargets(ctx, t, modIDs[0], true, []string{"Windows", "WindowsServer", "LinuxServer"})
		createVersionsForTargets(ctx, t, modIDs[2], false, []string{"Windows"})

		queryRequest := authRequest(`query ($modpackID: ModpackID!) {
			getModpackTargetSupport(modpackID: $modpackID) {
				target_name
			}
		}`, token)
		queryRequest.Var("modpackID", modpackID)

		var response struct {
			GetModpackTargetSupport []*struct {
				TargetName string `json:"target_name"`
			}
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &response))
		testza.AssertEqual(t, 3, len(response.GetModpackTargetSupport))

		targetNames := make(map[string]bool)
		for _, target := range response.GetModpackTargetSupport {
			targetNames[target.TargetName] = true
		}
		testza.AssertTrue(t, targetNames["Windows"])
		testza.AssertTrue(t, targetNames["WindowsServer"])
		testza.AssertTrue(t, targetNames["LinuxServer"])
	})

	t.Run("Missing target in optional dependency", func(t *testing.T) {
		modpackID := createTestModpack(ctx, t, client, token, userID, []string{modIDs[0], modIDs[3]})

		createVersionsForTargets(ctx, t, modIDs[0], true, []string{"Windows", "WindowsServer", "LinuxServer"})
		version3 := createVersionsForTargets(ctx, t, modIDs[3], false, []string{"Windows", "WindowsServer", "LinuxServer"})
		createVersionsForTargets(ctx, t, modIDs[1], true, []string{"Windows"})

		// Create optional dependency from mod3 to mod1
		_, err := db.From(ctx).VersionDependency.Create().
			SetVersionID(version3.ID).
			SetModID(modIDs[1]).
			SetCondition("1.0.0").
			SetOptional(true).
			Save(ctx)
		testza.AssertNoError(t, err)

		queryRequest := authRequest(`query ($modpackID: ModpackID!) {
			getModpackTargetSupport(modpackID: $modpackID) {
				target_name
			}
		}`, token)
		queryRequest.Var("modpackID", modpackID)

		var response struct {
			GetModpackTargetSupport []*struct {
				TargetName string `json:"target_name"`
			}
		}
		testza.AssertNoError(t, client.Run(ctx, queryRequest, &response))
		testza.AssertEqual(t, 3, len(response.GetModpackTargetSupport))

		targetNames := make(map[string]bool)
		for _, target := range response.GetModpackTargetSupport {
			targetNames[target.TargetName] = true
		}
		testza.AssertTrue(t, targetNames["Windows"])
		testza.AssertTrue(t, targetNames["WindowsServer"])
		testza.AssertTrue(t, targetNames["LinuxServer"])
	})

	t.Run("Missing target in non-optional dependency", func(t *testing.T) {
		modpackID := createTestModpack(ctx, t, client, token, userID, []string{modIDs[0], modIDs[3]})

		createVersionsForTargets(ctx, t, modIDs[0], true, []string{"Windows", "WindowsServer", "LinuxServer"})
		version3 := createVersionsForTargets(ctx, t, modIDs[3], false, []string{"Windows", "WindowsServer", "LinuxServer"})
		createVersionsForTargets(ctx, t, modIDs[1], true, []string{"Windows"})

		// Create non-optional dependency from mod3 to mod1
		_, err := db.From(ctx).VersionDependency.Create().
			SetVersionID(version3.ID).
			SetModID(modIDs[1]).
			SetCondition("1.0.0").
			SetOptional(false).
			Save(ctx)
		testza.AssertNoError(t, err)

		queryRequest := authRequest(`query ($modpackID: ModpackID!) {
			getModpackTargetSupport(modpackID: $modpackID) {
				target_name
			}
		}`, token)
		queryRequest.Var("modpackID", modpackID)

		var response struct {
			GetModpackTargetSupport []*struct {
				TargetName string `json:"target_name"`
			}
		}
		testza.AssertNotEqual(t, version3.Edges.VersionDependencies, nil)

		testza.AssertNoError(t, client.Run(ctx, queryRequest, &response))
		testza.AssertEqual(t, 1, len(response.GetModpackTargetSupport))
		testza.AssertEqual(t, "Windows", response.GetModpackTargetSupport[0].TargetName)
	})

}

func createTestModpack(ctx context.Context, t *testing.T, client *graphql.Client, token string, userID string, modIDs []string) string {
	mods := make([]struct {
		ModID             string `json:"mod_id"`
		VersionConstraint string `json:"version_constraint"`
	}, len(modIDs))

	for i, modID := range modIDs {
		mods[i] = struct {
			ModID             string `json:"mod_id"`
			VersionConstraint string `json:"version_constraint"`
		}{
			ModID:             modID,
			VersionConstraint: ">=0.0.0",
		}
	}

	createRequest := authRequest(`mutation ($name: String!, $shortDescription: String!, $fullDescription: String!, $targets: [String!]!, $mods: [ModpackModInput!]!) {
		createModpack(modpack: {
			name: $name,
			short_description: $shortDescription,
			full_description: $fullDescription,
			targets: $targets,
			mods: $mods
		}) {
			id
		}
	}`, token)

	createRequest.Var("name", "Test Modpack")
	createRequest.Var("shortDescription", "A test modpack")
	createRequest.Var("fullDescription", "A test modpack for testing target support")
	createRequest.Var("targets", []string{"Windows", "WindowsServer", "LinuxServer"})
	createRequest.Var("mods", mods)

	var response struct {
		CreateModpack generated.Modpack
	}
	testza.AssertNoError(t, client.Run(ctx, createRequest, &response))
	testza.AssertNotEqual(t, "", response.CreateModpack.ID)

	return response.CreateModpack.ID
}

func createVersionsForTargets(ctx context.Context, t *testing.T, modID string, requiredOnRemote bool, targets []string) *ent.Version {
	mod, err := db.From(ctx).Mod.Get(ctx, modID)
	testza.AssertNoError(t, err)

	version, err := db.From(ctx).Version.Create().
		SetModID(modID).
		SetVersion("1.0.0").
		SetGameVersion("1.0.0").
		SetModReference(mod.ModReference).
		SetRequiredOnRemote(requiredOnRemote).
		Save(ctx)
	testza.AssertNoError(t, err)

	// Create version targets
	for _, target := range targets {
		_, err := db.From(ctx).VersionTarget.Create().
			SetVersionID(version.ID).
			SetTargetName(target).
			SetHash("grumbus").
			SetSize(1000).
			Save(ctx)
		testza.AssertNoError(t, err)
	}

	return version
}
