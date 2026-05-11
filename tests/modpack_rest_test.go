package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MarvinJWendt/testza"
	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/config"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/nodes"
)

func init() {
	config.SetConfigDir("../")
	db.EnableDebug()
}

func setupEcho() *echo.Echo {
	e := echo.New()
	v1 := e.Group("/v1")
	nodes.RegisterModpackRoutes(v1.Group("/modpack"))
	nodes.RegisterModpacksRoutes(v1.Group("/modpacks"))
	return e
}

func TestModpacksRESTEndpoints(t *testing.T) {
	ctx, client, stop := setup()
	defer stop()

	e := setupEcho()

	token, _, err := makeUser(ctx)
	testza.AssertNoError(t, err)

	tags := seedTags(ctx, t, token, client)
	mods := seedMods(ctx, t, token, client, tags[0])
	modpacks := seedModpacks(ctx, t, token, client, tags, mods)

	t.Run("GET /v1/modpacks - List modpacks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/modpacks", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertTrue(t, response.Success)
		testza.AssertNil(t, response.Error)

		modpacksData, ok := response.Data.([]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, 2, len(modpacksData))
	})

	t.Run("GET /v1/modpacks - List modpacks with hidden=true", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/modpacks?hidden=true", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertTrue(t, response.Success)

		modpacksData, ok := response.Data.([]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, 3, len(modpacksData))
	})

	t.Run("GET /v1/modpacks - List modpacks with search", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/modpacks?search=Ultimate", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertTrue(t, response.Success)

		modpacksData, ok := response.Data.([]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, 1, len(modpacksData))
	})

	t.Run("GET /v1/modpacks - List modpacks with pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/modpacks?limit=1&offset=0&hidden=true", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertTrue(t, response.Success)

		modpacksData, ok := response.Data.([]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, 1, len(modpacksData))
	})

	t.Run("GET /v1/modpacks - List modpacks with ordering", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/modpacks?order_by=name&order=asc&hidden=true", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertTrue(t, response.Success)

		modpacksData, ok := response.Data.([]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, 3, len(modpacksData))

		firstModpack := modpacksData[0].(map[string]interface{})
		testza.AssertEqual(t, "Advanced Engineering", firstModpack["name"])
	})

	t.Run("GET /v1/modpack/:id - Get existing modpack", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/modpack/%s", modpacks[0]), nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertTrue(t, response.Success)
		testza.AssertNil(t, response.Error)

		modpackData, ok := response.Data.(map[string]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, "Ultimate Factory Pack", modpackData["name"])
		testza.AssertEqual(t, modpacks[0], modpackData["id"])

		tags, ok := modpackData["tags"].([]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertTrue(t, len(tags) > 0)

		mods, ok := modpackData["mods"].([]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, 3, len(mods))
	})

	t.Run("GET /v1/modpack/:id - Get non-existent modpack", func(t *testing.T) {
		nonExistentID := "non-existent-id"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/modpack/%s", nonExistentID), nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusNotFound, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertFalse(t, response.Success)
		testza.AssertNil(t, response.Data)

		errorData, ok := response.Error.(map[string]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, "modpack not found", errorData["message"])
		testza.AssertEqual(t, float64(400), errorData["code"])
	})

	t.Run("GET /v1/modpack/:id/releases/:version - Get existing release", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/modpack/%s/releases/1.0.0", modpacks[0]), nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertTrue(t, response.Success)
		testza.AssertNil(t, response.Error)

		releaseData, ok := response.Data.(map[string]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, "1.0.0", releaseData["version"])
		testza.AssertNotNil(t, releaseData["lockfile"])
	})

	t.Run("GET /v1/modpack/:id/releases/:version - Get non-existent release", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/modpack/%s/releases/999.0.0", modpacks[0]), nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusNotFound, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertFalse(t, response.Success)
		testza.AssertNil(t, response.Data)

		errorData, ok := response.Error.(map[string]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, "modpack release not found", errorData["message"])
		testza.AssertEqual(t, float64(401), errorData["code"])
	})

	t.Run("GET /v1/modpack/:id/releases/:version - Release for non-existent modpack", func(t *testing.T) {
		nonExistentID := "non-existent-modpack"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/modpack/%s/releases/1.0.0", nonExistentID), nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusNotFound, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertFalse(t, response.Success)
		testza.AssertNil(t, response.Data)

		errorData, ok := response.Error.(map[string]interface{})
		testza.AssertTrue(t, ok)
		testza.AssertEqual(t, "modpack release not found", errorData["message"])
	})
}

func TestModpackRESTValidation(t *testing.T) {
	ctx, _, stop := setup()
	defer stop()

	e := setupEcho()

	t.Run("GET /v1/modpacks - Invalid query parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/modpacks?limit=1000", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)

		req = httptest.NewRequest(http.MethodGet, "/v1/modpacks?order_by=invalid_field", nil)
		req = req.WithContext(ctx)
		rec = httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)

		req = httptest.NewRequest(http.MethodGet, "/v1/modpacks?order=invalid_order", nil)
		req = req.WithContext(ctx)
		rec = httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)
	})
}

func TestModpackRESTEdgeCases(t *testing.T) {
	ctx, _, stop := setup()
	defer stop()

	e := setupEcho()

	t.Run("GET /v1/modpacks - Empty search string", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/modpacks?search=", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)

		var response nodes.GenericResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		testza.AssertNoError(t, err)

		testza.AssertTrue(t, response.Success)
	})

	t.Run("GET /v1/modpacks - Zero limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/modpacks?limit=0", nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /v1/modpack/:id/releases/:version - Empty version", func(t *testing.T) {
		modpackID := "test-modpack-id"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/modpack/%s/releases/", modpackID), nil)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		testza.AssertEqual(t, http.StatusNotFound, rec.Code)
	})
}
