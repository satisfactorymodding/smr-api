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

	createRequest := authRequest(`mutation ($mod_reference: ModReference!) {
			createMod(mod: {
				name: "Rate Limit Test Mod",
				short_description: "Testing rate limiting functionality",
				full_description: "Lorem ipsum dolor sit amet",
				mod_reference: $mod_reference
			}) {
				id
				name
				mod_reference
			}
		}`, token)
	createRequest.Var("mod_reference", "newMod")

	var response struct {
		CreateMod generated.Mod
	}
	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, response.CreateMod)

	disclosureRequest := authRequest(`mutation ($mod_reference: ModReference!, $ai_use_disclosure: AIUseDisclosureInput!) {
		updateMod(mod: {
			name: "disclosure test",
			ai_use_disclosure: $ai_use_disclosure,
		}) {
			ai_use_disclosure {
				disclosure_type
				disclosure_string
			}
		}
	}`, token)
	disclosureString := "This mod uses AI for testing purposes"
	disclosureRequest.Var("mod_reference", "newMod")
	disclosureRequest.Var("ai_use_disclosure", generated.AIUseDisclosureInput{
		DisclosureType:   "ai_usage",
		DisclosureString: &disclosureString,
	})

	var createResponse struct {
		CreateMod generated.Mod
	}
	err = client.Run(ctx, disclosureRequest, &createResponse)
	testza.AssertNotNil(t, err)

	failedUpdate := authRequest(`mutation ($mod_reference: ModReference!, $ai_use_disclosure: AIUseDisclosureInput!) {
		updateMod(mod: {
			name: "disclosure test",
			ai_use_disclosure: $ai_use_disclosure,
		}) {
			ai_use_disclosure {
				disclosure_type
				disclosure_string
			}
		}
	}`, token)
	disclosureRequest.Var("mod_reference", "newMod")
	disclosureRequest.Var("ai_use_disclosure", nil)
	// Should not update since input is nil
	var failedResponse struct {
		CreateMod generated.Mod
	}
	err = client.Run(ctx, failedUpdate, &failedResponse)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, failedResponse.CreateMod.AiUseDisclosure.DisclosureType, "ai_usage")
	testza.AssertEqual(t, failedResponse.CreateMod.AiUseDisclosure.DisclosureString, disclosureString)
}
