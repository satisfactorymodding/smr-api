package tests

import (
	"testing"

	"github.com/MarvinJWendt/testza"

	"github.com/satisfactorymodding/smr-api/generated"
)

const getQuery = `query GetMods($offset: Int!, $limit: Int!, $search: String, $order: Order, $orderBy: ModFields, $tagIDs: [TagID!]) {
  getMods(
	filter: {limit: $limit, offset: $offset, search: $search, order: $order, order_by: $orderBy, tagIDs: $tagIDs}
  ) {
	count
	mods {
	  mod_reference
      tags {
        id
      }
	}
  }
}`

func TestGetModLimitOffset(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	tags := seedTags(ctx, t, token, client)
	seedMods(ctx, t, token, client, tags[0])

	getRequest := authRequest(getQuery, token)

	getRequest.Var("offset", "5")
	getRequest.Var("limit", "2")
	getRequest.Var("order", "asc")
	getRequest.Var("orderBy", "created_at")

	var getResponse struct {
		GetMods generated.GetMods
	}
	testza.AssertNoError(t, client.Run(ctx, getRequest, &getResponse))
	testza.AssertEqual(t, 11, getResponse.GetMods.Count)
	testza.AssertEqual(t, 2, len(getResponse.GetMods.Mods))
	testza.AssertEqual(t, "resource_overhaul", getResponse.GetMods.Mods[0].ModReference)
	testza.AssertEqual(t, "automated_defense", getResponse.GetMods.Mods[1].ModReference)
	testza.AssertEqual(t, 0, len(getResponse.GetMods.Mods[0].Tags))
	testza.AssertEqual(t, 1, len(getResponse.GetMods.Mods[1].Tags))
	testza.AssertEqual(t, tags[0], getResponse.GetMods.Mods[1].Tags[0].ID)
}
