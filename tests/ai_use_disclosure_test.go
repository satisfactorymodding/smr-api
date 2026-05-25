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
	modID := createResponse.CreateMod.ID

	// Empty string disclosure message is not allowed
	failedUpdateEmptyString := authRequest(`mutation ($id: ModID!, $ai_use_disclosure: AIUseDisclosureInput) {
		updateMod(
			modId: $id
			mod: {
				ai_use_disclosure: $ai_use_disclosure,
			}
		) {
			id
			ai_use_disclosure {
				disclosure_type
				message
			}
		}
	}`, token)
	failedUpdateEmptyString.Var("id", modID)
	emptyString := ""
	failedUpdateEmptyString.Var("ai_use_disclosure", generated.AIUseDisclosureInput{
		DisclosureType: generated.AIUseDisclosureTypeAiUsage,
		Message:        &emptyString,
	})

	var failedEmptyStringResponse struct {
		UpdateMod generated.Mod
	}
	err = client.Run(ctx, failedUpdateEmptyString, &failedEmptyStringResponse)
	testza.AssertNotNil(t, err)
	testza.AssertContains(t, err.Error(), "you need to input a disclosure message when disclosing AI usage")

	// Assigning a valid AI disclosure succeeds and updates the mod's disclosure information
	disclosureRequest := authRequest(`mutation ($id: ModID!, $ai_use_disclosure: AIUseDisclosureInput) {
		updateMod(
			modId: $id
			mod: {
				ai_use_disclosure: $ai_use_disclosure,
			}
		) {
			id
			ai_use_disclosure {
				disclosure_type
				message
			}
		}
	}`, token)
	disclosureMessage := "This mod uses AI for testing purposes"
	disclosureRequest.Var("id", modID)
	disclosureRequest.Var("ai_use_disclosure", generated.AIUseDisclosureInput{
		DisclosureType: generated.AIUseDisclosureTypeAiUsage,
		Message:        &disclosureMessage,
	})

	var updateResponse struct {
		UpdateMod generated.Mod
	}
	testza.AssertNoError(t, client.Run(ctx, disclosureRequest, &updateResponse))
	testza.AssertNotNil(t, updateResponse.UpdateMod.AiUseDisclosure.DisclosureType)
	testza.AssertEqual(t, generated.AIUseDisclosureTypeAiUsage, updateResponse.UpdateMod.AiUseDisclosure.DisclosureType)
	testza.AssertEqual(t, &disclosureMessage, updateResponse.UpdateMod.AiUseDisclosure.Message)

	// Trying to unassign the disclosure is not allowed
	failedUpdateUndisclosed := authRequest(`mutation ($id: ModID!, $ai_use_disclosure: AIUseDisclosureInput) {
		updateMod(
			modId: $id
			mod: {
				ai_use_disclosure: $ai_use_disclosure,
			}
		) {
			id
			ai_use_disclosure {
				disclosure_type
				message
			}
		}
	}`, token)
	failedUpdateUndisclosed.Var("id", modID)
	failedUpdateUndisclosed.Var("ai_use_disclosure", nil)

	var failedUpdateUndisclosedResponse struct {
		UpdateMod generated.Mod
	}
	err = client.Run(ctx, failedUpdateUndisclosed, &failedUpdateUndisclosedResponse)
	testza.AssertNotNil(t, err)
	testza.AssertContains(t, err.Error(), "this mod already has an AI use disclosure, and thus it cannot be cleared")
}
