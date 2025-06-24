package tests

import (
	"context"
	"testing"

	"github.com/MarvinJWendt/testza"
	"github.com/Masterminds/semver/v3"
	"github.com/machinebox/graphql"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/util"
)

func seedTags(ctx context.Context, t *testing.T, token string, client *graphql.Client) []string {
	tags := map[string]string{
		"hello": "hello N/A",
		"foo":   "foo N/A",
	}

	ids := make([]string, len(tags))
	i := 0
	for tag, desc := range tags {
		createRequest := authRequest(`mutation CreateTag($name: TagName!, $description: String!) {
		  createTag(tagName: $name, description: $description) {
			id
		  }
		}`, token)
		createRequest.Var("name", tag)
		createRequest.Var("description", desc)

		var createResponse struct {
			CreateTag generated.Tag
		}
		testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
		testza.AssertNotEqual(t, "", createResponse.CreateTag.ID)

		ids[i] = createResponse.CreateTag.ID
		i++
	}

	return ids
}

type testModpack struct {
	Name             string   `json:"name"`
	ShortDescription string   `json:"short_description"`
	FullDescription  string   `json:"full_description"`
	TagIDs           []string `json:"tagIDs"`
	Targets          []string `json:"targets"`
	Mods             []struct {
		ModID             string `json:"mod_id"`
		VersionConstraint string `json:"version_constraint"`
	} `json:"mods"`
	ParentID *string `json:"parent_id,omitempty"`
	Hidden   *bool   `json:"hidden,omitempty"`
}

func seedModpacks(ctx context.Context, t *testing.T, token string, client *graphql.Client, tagIDs []string, modIDs []string) []string {
	modpacks := []testModpack{
		{
			Name:             "Ultimate Factory Pack",
			ShortDescription: "The ultimate collection of factory enhancement mods",
			FullDescription:  "This modpack includes the best mods for optimizing and enhancing your factory operations, from automation to resource management.",
			TagIDs:           tagIDs,
			Targets:          []string{"Windows", "WindowsServer", "LinuxServer"},
			Mods: []struct {
				ModID             string `json:"mod_id"`
				VersionConstraint string `json:"version_constraint"`
			}{
				{ModID: modIDs[0], VersionConstraint: ">=1.0.0"},
				{ModID: modIDs[1], VersionConstraint: "^2.0.0"},
				{ModID: modIDs[2], VersionConstraint: "~1.0.0"},
			},
		},
		{
			Name:             "Beginner's Starter Pack",
			ShortDescription: "Perfect for new players getting started",
			FullDescription:  "A carefully curated selection of beginner-friendly mods that enhance the game experience without overwhelming new players.",
			TagIDs:           []string{tagIDs[0]},
			Targets:          []string{"Windows"},
			Mods: []struct {
				ModID             string `json:"mod_id"`
				VersionConstraint string `json:"version_constraint"`
			}{
				{ModID: modIDs[1], VersionConstraint: ">=2.0.0"},
				{ModID: modIDs[0], VersionConstraint: "^1.0.0"},
			},
		},
		{
			Name:             "Advanced Engineering",
			ShortDescription: "For experienced players seeking complex challenges",
			FullDescription:  "Advanced mods that add complexity and new engineering challenges for experienced players.",
			TagIDs:           []string{tagIDs[1]},
			Targets:          []string{"Windows", "WindowsServer", "LinuxServer"},
			Mods: []struct {
				ModID             string `json:"mod_id"`
				VersionConstraint string `json:"version_constraint"`
			}{
				{ModID: modIDs[1], VersionConstraint: ">=3.0.0"},
				{ModID: modIDs[2], VersionConstraint: "^1.0.0"},
				{ModID: modIDs[3], VersionConstraint: "~1.0.0"},
			},
			Hidden: func() *bool { b := true; return &b }(),
		},
	}

	ids := make([]string, len(modpacks))
	for i, pack := range modpacks {
		// Create the modpack itself
		createRequest := authRequest(`mutation ($modpack: NewModpack!) {
			createModpack(modpack: $modpack) {
				id
			}
		}`, token)
		createRequest.Var("modpack", pack)

		var createResponse struct {
			CreateModpack generated.Modpack
		}
		testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
		testza.AssertNotEqual(t, "", createResponse.CreateModpack.ID)

		ids[i] = createResponse.CreateModpack.ID

		// Create a release for each modpack
		releaseRequest := authRequest(`mutation ($modpackID: ModpackID!, $release: NewModpackRelease!) {
			createModpackRelease(modpackID: $modpackID, release: $release) {
				id
			}
		}`, token)
		releaseRequest.Var("modpackID", createResponse.CreateModpack.ID)
		releaseRequest.Var("release", struct {
			Version   string `json:"version"`
			Changelog string `json:"changelog"`
		}{
			Version:   "1.0.0",
			Changelog: "Hello World",
		})

		var releaseaResponse struct {
			CreateModpackRelease generated.ModpackRelease
		}
		testza.AssertNoError(t, client.Run(ctx, releaseRequest, &releaseaResponse))
		testza.AssertNotEqual(t, "", releaseaResponse.CreateModpackRelease.ID)
	}

	return ids
}

type testVersion struct {
	Version      *semver.Version `json:"version"`
	Dependencies map[string]string
}

type testMod struct {
	Name             string   `json:"name"`
	ShortDescription string   `json:"short_description"`
	FullDescription  string   `json:"full_description"`
	ModReference     string   `json:"mod_reference"`
	TagIDs           []string `json:"tagIDs"`

	versions []testVersion
}

func seedMods(ctx context.Context, t *testing.T, token string, client *graphql.Client, tagID string) []string {
	mods := []testMod{
		{
			Name:             "Advanced Robotics",
			ShortDescription: "Enhances robot efficiency and adds new automation features.",
			ModReference:     "advanced_robotics",
			versions: []testVersion{
				{
					Version: semver.New(1, 0, 0, "", ""),
					Dependencies: map[string]string{
						"mega_factory": ">0.0.0",
					},
				},
				{
					Version: semver.New(2, 0, 0, "", ""),
					Dependencies: map[string]string{
						"mega_factory": ">0.0.0",
					},
				},
				{
					Version: semver.New(3, 0, 0, "", ""),
					Dependencies: map[string]string{
						"mega_factory": ">0.0.0",
					},
				},
			},
		},
		{
			Name:             "Eco-Friendly Power",
			ShortDescription: "Introduces sustainable energy sources and eco-friendly power management.",
			ModReference:     "eco_friendly_power",
			versions: []testVersion{
				{
					Version: semver.New(1, 0, 0, "", ""),
					Dependencies: map[string]string{
						"quantum_transport": ">0.0.0",
					},
				},
				{
					Version: semver.New(2, 0, 0, "", ""),
					Dependencies: map[string]string{
						"quantum_transport": ">0.0.0",
					},
				},
				{
					Version: semver.New(3, 0, 0, "", ""),
					Dependencies: map[string]string{
						"quantum_transport": ">0.0.0",
					},
				},
			},
		},
		{
			Name:             "Quantum Transport",
			ShortDescription: "Allows instantaneous item transport using quantum entanglement.",
			ModReference:     "quantum_transport",
			versions: []testVersion{
				{
					Version: semver.New(1, 0, 0, "", ""),
				},
			},
		},
		{
			Name:             "Mega Factory",
			ShortDescription: "Expands factory building limits and adds new large-scale production tools.",
			ModReference:     "mega_factory",
			versions: []testVersion{
				{
					Version: semver.New(1, 0, 0, "", ""),
				},
			},
		},
		{
			Name:             "Resource Overhaul",
			ShortDescription: "Revamps resource extraction and processing for more efficiency.",
			ModReference:     "resource_overhaul",
		},
		{
			Name:             "Automated Defense",
			ShortDescription: "Adds advanced automated defense systems to protect your factory.",
			ModReference:     "automated_defense",
			TagIDs:           []string{tagID},
		},
		{
			Name:             "AI Assistant",
			ShortDescription: "Introduces an AI assistant to help manage and optimize your factory.",
			ModReference:     "ai_assistant",
			TagIDs:           []string{tagID},
		},
		{
			Name:             "Fusion Reactors",
			ShortDescription: "Adds fusion reactors as a high-efficiency power source.",
			ModReference:     "fusion_reactors",
			TagIDs:           []string{tagID},
		},
		{
			Name:             "Modular Production",
			ShortDescription: "Allows modular production units for flexible factory layouts.",
			ModReference:     "modular_production",
			TagIDs:           []string{tagID},
		},
		{
			Name:             "Nanotech Manufacturing",
			ShortDescription: "Incorporates nanotechnology for ultra-precise manufacturing processes.",
			ModReference:     "nanotech_manufacturing",
			TagIDs:           []string{tagID},
		},
	}

	util.ModsPer24h = len(mods)

	ids := make([]string, len(mods))
	for i, mod := range mods {
		mod.FullDescription = "N/A"

		createRequest := authRequest(`mutation CreateMod($mod: NewMod!) {
		  createMod(mod: $mod) {
			id
		  }
		}`, token)
		createRequest.Var("mod", mod)

		var createResponse struct {
			CreateMod generated.Mod
		}
		testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
		testza.AssertNotEqual(t, "", createResponse.CreateMod.ID)

		ids[i] = createResponse.CreateMod.ID

		for _, version := range mod.versions {
			v := db.From(ctx).Version.Create().
				SetVersion(version.Version.String()).
				SetVersionMajor(int(version.Version.Major())).
				SetVersionMinor(int(version.Version.Minor())).
				SetVersionPatch(int(version.Version.Patch())).
				SetGameVersion(">=0").
				SetModReference(mod.ModReference).
				SetModID(createResponse.CreateMod.ID).
				SetApproved(true).
				SaveX(ctx)

			for _, s := range []string{"Windows", "WindowsServer", "LinuxServer"} {
				db.From(ctx).VersionTarget.Create().SetVersion(v).SetTargetName(s).SaveX(ctx)
			}
		}
	}

	return ids
}
