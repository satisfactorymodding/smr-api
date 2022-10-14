package nodes

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
)

// @Summary Retrieve a list of Mods
// @Tags Mods
// @Description Retrieve a list of mods
// @Accept  json
// @Produce  json
// @Param limit query int false "How many mods to return"
// @Param offset query int false "Offset for list of mods to return"
// @Param order_by query string false "Order by field" Enums(created_at, updated_at, name, views, downloads, hotness, popularity, last_version_date)
// @Param order query string false "Order of results" Enums(asc, desc)
// @Param search query string false "Search string"
// @Success 200
// @Router /mods [get]
func getMods(c echo.Context) (interface{}, *ErrorResponse) {
	limit := util.GetIntRange(c, "limit", 1, 100, 25)
	offset := util.GetIntRange(c, "offset", 0, 9999999, 0)
	orderBy := util.OneOf(c, "order_by", []string{"created_at", "updated_at", "name", "views", "downloads", "hotness", "popularity", "last_version_date"}, "created_at")
	order := util.OneOf(c, "order", []string{"asc", "desc"}, "desc")
	search := c.QueryParam("search")

	mods := postgres.GetMods(c.Request().Context(), limit, offset, orderBy, order, search, false)

	converted := make([]*Mod, len(mods))
	for k, v := range mods {
		converted[k] = ModToMod(&v, true)
	}

	return converted, nil
}

// @Summary Retrieve a count of Mods
// @Tags Mods
// @Description Retrieve a count of Mods
// @Accept  json
// @Produce  json
// @Param search query string false "Search string"
// @Success 200
// @Router /mods/count [get]
func getModCount(c echo.Context) (interface{}, *ErrorResponse) {
	search := c.QueryParam("search")
	return postgres.GetModCount(c.Request().Context(), search, false), nil
}

// @Summary Retrieve a Mod
// @Tags Mod
// @Description Retrieve a mod by mod ID
// @Accept  json
// @Produce  json
// @Param modId path string true "Mod ID"
// @Success 200
// @Router /mod/{modId} [get]
func getMod(c echo.Context) (interface{}, *ErrorResponse) {
	modID := c.Param("modId")

	mod := postgres.GetModByID(c.Request().Context(), modID)

	if mod == nil {
		return nil, &ErrorModNotFound
	}

	if _, ok := c.QueryParams()["view"]; ok {
		if redis.CanIncrement(c.RealIP(), "view", "mod:"+modID, time.Hour*4) {
			postgres.IncrementModViews(c.Request().Context(), mod)
		}
	}

	return ModToMod(mod, false), nil
}

// @Summary Retrieve a list of Mods by ID
// @Tags Mods
// @Description Retrieve a list of mods by mod IDs
// @Accept  json
// @Produce  json
// @Param modIds path string true "Mod IDs"
// @Success 200
// @Router /mods/{modIds} [get]
func getModsByIds(c echo.Context) (interface{}, *ErrorResponse) {
	modID := c.Param("modIds")
	modIDSplit := strings.Split(modID, ",")

	// TODO limit amount of users requestable

	mods := postgres.GetModsByID(c.Request().Context(), modIDSplit)

	if mods == nil {
		return nil, &ErrorModNotFound
	}

	converted := make([]*Mod, len(mods))
	for k, v := range mods {
		converted[k] = ModToMod(&v, true)
	}

	return converted, nil
}

// @Summary Retrieve a list of latest versions for a mod
// @Tags Mod
// @Description Retrieve a list of latest versions for a mod based on mod id
// @Accept  json
// @Produce  json
// @Param modId path string true "Mod ID"
// @Success 200
// @Router /mod/{modId}/latest-versions [get]
func getModLatestVersions(c echo.Context) (interface{}, *ErrorResponse) {
	modID := c.Param("modId")

	versions := postgres.GetModsLatestVersions(c.Request().Context(), []string{modID}, false)

	if versions == nil {
		return nil, &ErrorVersionNotFound
	}

	result := make(map[string]*Version)

	for _, v := range *versions {
		result[v.Stability] = VersionToVersion(&v)
	}

	return result, nil
}

// @Summary Retrieve a list of latest versions for mods
// @Tags Mods
// @Description Retrieve a list of latest versions for mods based on mod id
// @Accept  json
// @Produce  json
// @Param modIds path string true "Mod IDs"
// @Success 200
// @Router /mods/{modIds}/latest-versions [get]
func getModsLatestVersions(c echo.Context) (interface{}, *ErrorResponse) {
	modID := c.Param("modIds")
	modIDSplit := strings.Split(modID, ",")

	// TODO limit amount of mods requestable

	versions := postgres.GetModsLatestVersions(c.Request().Context(), modIDSplit, false)

	if versions == nil {
		return nil, &ErrorVersionNotFound
	}

	result := make(map[string]map[string]*Version)

	for _, v := range *versions {
		if _, ok := result[v.ModID]; !ok {
			result[v.ModID] = make(map[string]*Version)
		}
		result[v.ModID][v.Stability] = VersionToVersion(&v)
	}

	return result, nil
}

// @Summary Retrieve a Mod Versions
// @Tags Mod
// @Description Retrieve a mod versions by mod ID
// @Accept  json
// @Produce  json
// @Param limit query int false "How many versions to return"
// @Param offset query int false "Offset for list of versions to return"
// @Param order_by query string false "Order by field" Enums(created_at, updated_at)
// @Param order query string false "Order of results" Enums(asc, desc)
// @Param modId path string true "Mod ID"
// @Success 200
// @Router /mod/{modId}/versions [get]
func getModVersions(c echo.Context) (interface{}, *ErrorResponse) {
	limit := util.GetIntRange(c, "limit", 1, 100, 25)
	offset := util.GetIntRange(c, "offset", 0, 9999999, 0)
	orderBy := util.OneOf(c, "order_by", []string{"created_at", "updated_at"}, "created_at")
	order := util.OneOf(c, "order", []string{"asc", "desc"}, "desc")

	modID := c.Param("modId")

	mod := postgres.GetModByID(c.Request().Context(), modID)

	if mod == nil {
		return nil, &ErrorModNotFound
	}

	versions := postgres.GetModVersions(c.Request().Context(), mod.ID, limit, offset, orderBy, order, false)

	converted := make([]*Version, len(versions))
	for k, v := range versions {
		converted[k] = VersionToVersion(&v)
	}

	return converted, nil
}

// @Summary Retrieve a Mod Authors
// @Tags Mod
// @Description Retrieve a mod authors by mod ID
// @Accept  json
// @Produce  json
// @Param modId path string true "Mod ID"
// @Success 200
// @Router /mod/{modId}/authors [get]
func getModAuthors(c echo.Context) (interface{}, *ErrorResponse) {
	modID := c.Param("modId")

	mod := postgres.GetModByID(c.Request().Context(), modID)

	if mod == nil {
		return nil, &ErrorModNotFound
	}

	authors := postgres.GetModAuthors(c.Request().Context(), mod.ID)

	converted := make([]*ModUser, len(authors))
	for k, v := range authors {
		converted[k] = ModUserToModUser(&v)
	}

	return converted, nil
}

// @Summary Retrieve a Mod Version
// @Tags Mod
// @Description Retrieve a mod version by mod ID and version ID
// @Accept  json
// @Produce  json
// @Param modId path string true "Mod ID"
// @Param versionId path string true "Version ID"
// @Success 200
// @Router /mod/{modId}/versions/{versionId} [get]
func getModVersion(c echo.Context) (interface{}, *ErrorResponse) {
	modID := c.Param("modId")
	versionID := c.Param("versionId")

	mod := postgres.GetModByID(c.Request().Context(), modID)

	if mod == nil {
		return nil, &ErrorModNotFound
	}

	version := postgres.GetModVersion(c.Request().Context(), mod.ID, versionID)

	if version == nil {
		return nil, &ErrorVersionNotFound
	}

	return VersionToVersion(version), nil
}

// @Summary Download a Mod Version
// @Tags Mod
// @Description Download a mod version by mod ID and version ID
// @Accept  json
// @Produce  json
// @Param modId path string true "Mod ID"
// @Param versionId path string true "Version ID"
// @Success 200
// @Router /mod/{modId}/versions/{versionId}/download [get]
func downloadModVersion(c echo.Context) error {
	modID := c.Param("modId")
	versionID := c.Param("versionId")

	mod := postgres.GetModByID(c.Request().Context(), modID)

	if mod == nil {
		return c.String(404, "mod not found, modID:"+modID)
	}

	version := postgres.GetModVersion(c.Request().Context(), mod.ID, versionID)

	if version == nil {
		return c.String(404, "version not found, modID:"+modID+" versionID:"+versionID)
	}

	if redis.CanIncrement(c.RealIP(), "download", "version:"+versionID, time.Hour*4) {
		postgres.IncrementVersionDownloads(c.Request().Context(), version)
	}

	return c.Redirect(302, storage.GenerateDownloadLink(version.Key))
}

// @Summary Download a Mod Version by Platform
// @Tags Mod
// @Description Download a mod version by mod ID and version ID and Platform
// @Accept  json
// @Produce  json
// @Param modId path string true "Mod ID"
// @Param versionId path string true "Version ID"
// @Param versionId path string true "Platform"
// @Success 200
// @Router /mod/{modId}/versions/{versionId}/{platform}/download [get]
func downloadModVersionArch(c echo.Context) error {
	modID := c.Param("modId")
	versionID := c.Param("versionId")
	platform := c.Param("platform")

	mod := postgres.GetModByID(c.Request().Context(), modID)

	if mod == nil {
		return c.String(404, "mod not found, modID:"+modID)
	}

	version := postgres.GetModVersion(c.Request().Context(), mod.ID, versionID)

	if version == nil {
		return c.String(404, "version not found, modID:"+modID+" versionID:"+versionID)
	}

	arch := postgres.GetModArchByPlatform(c.Request().Context(), versionID, platform)

	if arch == nil {
		return c.String(404, "platform not found, modID:"+modID+" versionID:"+versionID+" platform:"+platform)
	}

	if redis.CanIncrement(c.RealIP(), "download", "version:"+versionID, time.Hour*4) {
		postgres.IncrementVersionDownloads(c.Request().Context(), version)
	}

	return c.Redirect(302, storage.GenerateDownloadLink(arch.Key))
}
