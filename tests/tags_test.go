package tests

import (
	"strconv"
	"testing"

	"github.com/MarvinJWendt/testza"

	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/migrations"
)

func init() {
	migrations.SetMigrationDir("../migrations")
	config.SetConfigDir("../")
	db.EnableDebug()
}

func TestTags(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	// Run Twice to detect any cache issues
	for i := range 2 {
		t.Run("Loop"+strconv.Itoa(i), func(t *testing.T) {
			var objID string

			var firstTag string
			var secondTag string

			t.Run("Create", func(t *testing.T) {
				createRequest := authRequest(`mutation CreateTag($tagName: TagName!, $description: String!) {
				  obj: createTag(tagName: $tagName, description: $description) {
					id
					name
					description
				  }
				}`, token)
				createRequest.Var("tagName", "Foo"+strconv.Itoa(i))
				createRequest.Var("description", "Lorem Ipsum")

				var createResponse struct {
					Obj generated.Tag
				}
				testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))
				testza.AssertNotEqual(t, "", createResponse.Obj.ID)
				testza.AssertEqual(t, "Foo"+strconv.Itoa(i), createResponse.Obj.Name)
				testza.AssertEqual(t, "Lorem Ipsum", createResponse.Obj.Description)

				objID = createResponse.Obj.ID
			})

			t.Run("Create Many", func(t *testing.T) {
				createRequest := authRequest(`mutation CreateMultipleTags($tagNames: [NewTag!]!) {
				  obj: createMultipleTags(tagNames: $tagNames) {
					id
					name
					description
				  }
				}`, token)
				createRequest.Var("tagNames", []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
				}{
					{
						Name:        "One" + strconv.Itoa(i),
						Description: "First Tag",
					},
					{
						Name:        "Two" + strconv.Itoa(i),
						Description: "Second Tag",
					},
				})

				var createResponse struct {
					Obj []generated.Tag
				}
				testza.AssertNoError(t, client.Run(ctx, createRequest, &createResponse))

				testza.AssertNotEqual(t, "", createResponse.Obj[0].ID)
				testza.AssertEqual(t, "One"+strconv.Itoa(i), createResponse.Obj[0].Name)
				testza.AssertEqual(t, "First Tag", createResponse.Obj[0].Description)

				testza.AssertNotEqual(t, "", createResponse.Obj[1].ID)
				testza.AssertEqual(t, "Two"+strconv.Itoa(i), createResponse.Obj[1].Name)
				testza.AssertEqual(t, "Second Tag", createResponse.Obj[1].Description)

				firstTag = createResponse.Obj[0].ID
				secondTag = createResponse.Obj[1].ID
			})

			t.Run("Query One", func(t *testing.T) {
				queryRequest := authRequest(`query GetTag($tagId: TagID!) {
				  obj: getTag(tagID: $tagId) {
					id
					name
					description
				  }
				}`, token)
				queryRequest.Var("tagId", objID)

				var queryResponse struct {
					Obj generated.Tag
				}
				testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
				testza.AssertEqual(t, objID, queryResponse.Obj.ID)
				testza.AssertEqual(t, "Foo"+strconv.Itoa(i), queryResponse.Obj.Name)
				testza.AssertEqual(t, "Lorem Ipsum", queryResponse.Obj.Description)
			})

			t.Run("Update", func(t *testing.T) {
				updateRequest := authRequest(`mutation UpdateTag($tagId: TagID!, $newName: TagName!, $description: String!) {
				  obj: updateTag(tagID: $tagId, NewName: $newName, description: $description) {
					id
					name
					description
				  }
				}`, token)
				updateRequest.Var("tagId", objID)
				updateRequest.Var("newName", "Meow"+strconv.Itoa(i))
				updateRequest.Var("description", "I'm a teapot")

				var updateResponse struct {
					Obj generated.Tag
				}
				testza.AssertNoError(t, client.Run(ctx, updateRequest, &updateResponse))
				testza.AssertEqual(t, objID, updateResponse.Obj.ID)
				testza.AssertEqual(t, "Meow"+strconv.Itoa(i), updateResponse.Obj.Name)
				testza.AssertEqual(t, "I'm a teapot", updateResponse.Obj.Description)
			})

			t.Run("Query Many", func(t *testing.T) {
				queryRequest := authRequest(`query {
				  obj: getTags(filter: {order: asc}) {
					id
					name
					description
				  }
				}`, token)

				var queryResponse struct {
					Obj []generated.Tag
				}
				testza.AssertNoError(t, client.Run(ctx, queryRequest, &queryResponse))
				testza.AssertEqual(t, 3, len(queryResponse.Obj))
				testza.AssertEqual(t, objID, queryResponse.Obj[0].ID)
				testza.AssertEqual(t, "Meow"+strconv.Itoa(i), queryResponse.Obj[0].Name)
				testza.AssertEqual(t, "I'm a teapot", queryResponse.Obj[0].Description)
			})

			t.Run("Delete", func(t *testing.T) {
				for _, s := range []string{objID, firstTag, secondTag} {
					deleteRequest := authRequest(`mutation DeleteTag($tagId: TagID!) {
					  response: deleteTag(tagID: $tagId)
					}`, token)
					deleteRequest.Var("tagId", s)

					var deleteResponse struct {
						Response bool
					}
					testza.AssertNoError(t, client.Run(ctx, deleteRequest, &deleteResponse))
					testza.AssertTrue(t, deleteResponse.Response)
				}
			})
		})
	}
}
