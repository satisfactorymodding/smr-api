package tests

import (
	"context"
	"testing"

	"github.com/MarvinJWendt/testza"
	"github.com/machinebox/graphql"

	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/util"
)

func seedTags(ctx context.Context, t *testing.T, token string, client *graphql.Client) []string {
	tags := []string{
		"hello",
		"foo",
	}

	ids := make([]string, len(tags))
	for i, tag := range tags {
		createRequest := authRequest(`mutation CreateTag($name: TagName!) {
		  createTag(tagName: $name, description: "N/A") {
			id
		  }
		}`, token)
		createRequest.Var("name", tag)

		var createResponse struct {
			CreateTag generated.Tag
		}
		testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
		testza.AssertNotEqual(t, "", createResponse.CreateTag.ID)

		ids[i] = createResponse.CreateTag.ID
	}

	return ids
}

type testMod struct {
	Name             string   `json:"name"`
	ShortDescription string   `json:"short_description"`
	FullDescription  string   `json:"full_description"`
	ModReference     string   `json:"mod_reference"`
	TagIDs           []string `json:"tagIDs"`
}

func seedMods(ctx context.Context, t *testing.T, token string, client *graphql.Client, tagID string) []string {
	mods := []testMod{
		{
			Name:             "Advanced Robotics",
			ShortDescription: "Enhances robot efficiency and adds new automation features.",
			ModReference:     "advanced_robotics",
		},
		{
			Name:             "Eco-Friendly Power",
			ShortDescription: "Introduces sustainable energy sources and eco-friendly power management.",
			ModReference:     "eco_friendly_power",
		},
		{
			Name:             "Quantum Transport",
			ShortDescription: "Allows instantaneous item transport using quantum entanglement.",
			ModReference:     "quantum_transport",
		},
		{
			Name:             "Mega Factory",
			ShortDescription: "Expands factory building limits and adds new large-scale production tools.",
			ModReference:     "mega_factory",
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
	}

	return ids
}
