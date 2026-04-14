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

func TestAiDisclosure(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	// Creating a mod without an AI disclosure succeeds and there is no disclosure information
	createRequest := authRequest(`mutation ($mod_reference: ModReference!) {
			createMod(mod: {
				name: "AI Disclosure Mod",
				short_description: "Testing AI disclosure functionality",
				full_description: "Lorem ipsum dolor sit amet",
				mod_reference: $mod_reference
			}) {
				id
				name
				mod_reference
			}
		}`, token)
	createRequest.Var("mod_reference", "aiDisclosureMod")

	var createResponse struct {
		CreateMod generated.Mod
	}
	testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
	testza.AssertNotNil(t, createResponse.CreateMod)
	testza.AssertNil(t, createResponse.CreateMod.AiUseDisclosure)
	modId := createResponse.CreateMod.ID

	// Assigning an AI disclosure succeeds and updates the mod's disclosure information
	disclosureRequest := authRequest(`mutation ($id: ModID!, $ai_use_disclosure: AIUseDisclosureInput!) {
		updateMod(
			modId: $id
			mod: {
				ai_use_disclosure: $ai_use_disclosure,
			}
		) {
			id
			ai_use_disclosure {
				disclosure_type
				disclosure_string
			}
		}
	}`, token)
	disclosureString := "This mod uses AI for testing purposes"
	disclosureRequest.Var("id", modId)
	disclosureRequest.Var("ai_use_disclosure", generated.AIUseDisclosureInput{
		DisclosureType:   "ai_usage",
		DisclosureString: &disclosureString,
	})

	var updateResponse struct {
		UpdateMod generated.Mod
	}
	testza.AssertNoError(t, client.Run(ctx, disclosureRequest, &updateResponse))
	testza.AssertNotNil(t, updateResponse.UpdateMod.AiUseDisclosure.DisclosureType)
	testza.AssertEqual(t, generated.AIUseDisclosureTypeAiUsage, updateResponse.UpdateMod.AiUseDisclosure.DisclosureType)
	testza.AssertEqual(t, &disclosureString, updateResponse.UpdateMod.AiUseDisclosure.DisclosureString)

	// Trying to set to nil is not allowed
	failedUpdateNil := authRequest(`mutation ($id: ModID!, $ai_use_disclosure: AIUseDisclosureInput!) {
		updateMod(
			modId: $id
			mod: {
				ai_use_disclosure: $ai_use_disclosure,
			}
		) {
			id
			ai_use_disclosure {
				disclosure_type
				disclosure_string
			}
		}
	}`, token)
	failedUpdateNil.Var("id", modId)
	failedUpdateNil.Var("ai_use_disclosure", nil)

	var failedNilResponse struct {
		UpdateMod generated.Mod
	}
	err = client.Run(ctx, failedUpdateNil, &failedNilResponse)
	testza.AssertNotNil(t, err)
	testza.AssertContains(t, err.Error(), "cannot be null")
}
