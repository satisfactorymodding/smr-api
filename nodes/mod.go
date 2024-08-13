package nodes

import (
	"log/slog"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Vilsol/slox"
	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	mod2 "github.com/satisfactorymodding/smr-api/generated/ent/mod"
	version2 "github.com/satisfactorymodding/smr-api/generated/ent/version"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiontarget"
	"github.com/satisfactorymodding/smr-api/models"
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

	modFilter := models.DefaultModFilter()
	modFilter.Limit = &limit
	modFilter.Offset = &offset

	orderByGen := generated.ModFields(orderBy)
	modFilter.OrderBy = &orderByGen

	orderGen := generated.Order(order)
	modFilter.Order = &orderGen
	modFilter.Search = &search

	query := db.From(c.Request().Context()).Mod.Query()
	query = db.ConvertModFilter(query, modFilter, false, false)

	mods, err := query.All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed retrieving mods", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	return (*conv.ModImpl)(nil).ConvertSlice(mods), nil
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

	modFilter := models.DefaultModFilter()
	modFilter.Search = &search

	query := db.From(c.Request().Context()).Mod.Query()
	query = db.ConvertModFilter(query, modFilter, true, false)

	count, err := query.Count(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed retrieving mod count", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	return count, nil
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

	mod, err := db.From(c.Request().Context()).Mod.Query().
		WithTags().
		Where(mod2.ID(modID)).
		First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mod", slog.Any("err", err))
		return nil, &ErrorModNotFound
	}

	if mod == nil {
		return nil, &ErrorModNotFound
	}

	if _, ok := c.QueryParams()["view"]; ok {
		if redis.CanIncrement(c.RealIP(), "view", "mod:"+modID, time.Hour*4) {
			_ = mod.Update().AddViews(1).Exec(c.Request().Context())
		}
	}

	return (*conv.ModImpl)(nil).Convert(mod), nil
}

// @Summary Retrieve a list of Mods by ID
// @Tags Mods
// @Description Retrieve a list of mods by mod IDs
// @Accept  json
// @Produce  json
// @Param modIds path string true "Mod IDs"
// @Success 200
// @Router /mods/{modIds} [get]
func getModsByIDs(c echo.Context) (interface{}, *ErrorResponse) {
	modID := c.Param("modIds")
	modIDSplit := strings.Split(modID, ",")

	// TODO limit amount of users requestable

	mods, err := db.From(c.Request().Context()).Mod.Query().
		WithTags().
		Where(mod2.IDIn(modIDSplit...)).
		All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mods", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	if mods == nil {
		return nil, &ErrorModNotFound
	}

	return (*conv.ModImpl)(nil).ConvertSlice(mods), nil
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

	versions, err := db.From(c.Request().Context()).Version.Query().
		WithTargets().
		Modify(func(s *sql.Selector) {
			s.SelectExpr(sql.ExprP("distinct on (mod_id, stability) *"))
		}).
		Where(version2.Approved(true), version2.Denied(false), version2.ModID(modID)).
		Order(version2.ByStability(sql.OrderDesc()), version2.ByCreatedAt(sql.OrderDesc())).
		All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching versions", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	if versions == nil {
		return nil, &ErrorVersionNotFound
	}

	result := make(map[string]*generated.Version)

	for _, v := range versions {
		result[string(v.Stability)] = (*conv.VersionImpl)(nil).Convert(v)
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

	versions, err := db.From(c.Request().Context()).Version.Query().
		WithTargets().
		Modify(func(s *sql.Selector) {
			s.SelectExpr(sql.ExprP("distinct on (mod_id, stability) *"))
		}).
		Where(version2.Approved(true), version2.Denied(false), version2.ModIDIn(modIDSplit...)).
		Order(version2.ByStability(sql.OrderDesc()), version2.ByCreatedAt(sql.OrderDesc())).
		All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching versions", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	if versions == nil {
		return nil, &ErrorVersionNotFound
	}

	result := make(map[string]map[string]*generated.Version)

	for _, v := range versions {
		if _, ok := result[v.ModID]; !ok {
			result[v.ModID] = make(map[string]*generated.Version)
		}
		result[v.ModID][string(v.Stability)] = (*conv.VersionImpl)(nil).Convert(v)
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

	mod, err := db.From(c.Request().Context()).Mod.Query().
		WithTags().
		Where(mod2.ID(modID)).
		First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mod", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	if mod == nil {
		return nil, &ErrorModNotFound
	}

	versions, err := mod.QueryVersions().
		WithDependencies().
		WithTargets().
		Limit(limit).
		Offset(offset).
		Order(sql.OrderByField(orderBy, db.OrderToOrder(order)).ToFunc()).
		Where(version2.Approved(true), version2.Denied(false)).
		All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching versions", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	return (*conv.VersionImpl)(nil).ConvertSlice(versions), nil
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

	mod, err := db.From(c.Request().Context()).Mod.Query().
		WithTags().
		Where(mod2.ID(modID)).
		First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mod", slog.Any("err", err))
		return nil, &ErrorModNotFound
	}

	if mod == nil {
		return nil, &ErrorModNotFound
	}

	authors, err := mod.QueryUserMods().All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching authors", slog.Any("err", err))
		return nil, &ErrorModNotFound
	}

	return (*conv.UserModImpl)(nil).ConvertSlice(authors), nil
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

	mod, err := db.From(c.Request().Context()).Mod.Query().
		WithTags().
		Where(mod2.ID(modID)).
		First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mod", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	if mod == nil {
		return nil, &ErrorModNotFound
	}

	version, err := mod.QueryVersions().Where(version2.ID(versionID)).First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching version", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	if version == nil {
		return nil, &ErrorVersionNotFound
	}

	return (*conv.VersionImpl)(nil).Convert(version), nil
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

	mod, err := db.From(c.Request().Context()).Mod.Query().
		WithTags().
		Where(mod2.ID(modID)).
		First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mod", slog.Any("err", err))
		return c.String(404, "mod not found, modID:"+modID)
	}

	if mod == nil {
		return c.String(404, "mod not found, modID:"+modID)
	}

	version, err := mod.QueryVersions().Where(version2.ID(versionID)).First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching version", slog.Any("err", err))
		return c.String(404, "version not found, modID:"+modID+" versionID:"+versionID)
	}

	if version == nil {
		return c.String(404, "version not found, modID:"+modID+" versionID:"+versionID)
	}

	if c.Request().Method == echo.GET &&
		redis.CanIncrement(c.RealIP(), "download", "version:"+versionID, time.Hour*4) {
		_ = version.Update().AddDownloads(1).Exec(c.Request().Context())
	}

	return c.Redirect(302, storage.GenerateDownloadLink(version.Key))
}

// @Summary Download a Mod Version by TargetName
// @Tags Mod
// @Description Download a mod version by mod ID and version ID and TargetName
// @Accept  json
// @Produce  json
// @Param modId path string true "Mod ID"
// @Param versionId path string true "Version ID"
// @Param target path string true "TargetName"
// @Success 200
// @Router /mod/{modId}/versions/{versionId}/{target}/download [get]
func downloadModVersionTarget(c echo.Context) error {
	modID := c.Param("modId")
	versionID := c.Param("versionId")
	target := c.Param("target")

	mod, err := db.From(c.Request().Context()).Mod.Query().
		WithTags().
		Where(mod2.ID(modID)).
		First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mod", slog.Any("err", err))
		return c.String(404, "mod not found, modID:"+modID)
	}

	if mod == nil {
		return c.String(404, "mod not found, modID:"+modID)
	}

	version, err := mod.QueryVersions().Where(version2.ID(versionID)).First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching version", slog.Any("err", err))
		return err
	}

	if version == nil {
		return c.String(404, "version not found, modID:"+modID+" versionID:"+versionID)
	}

	versionTarget, err := version.QueryTargets().Where(versiontarget.TargetName(target)).First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching target", slog.Any("err", err))
		return err
	}

	if versionTarget == nil {
		return c.String(404, "target not found, modID:"+modID+" versionID:"+versionID+" target:"+target)
	}

	if c.Request().Method == echo.GET &&
		redis.CanIncrement(c.RealIP(), "download", "version:"+versionID, time.Hour*4) {
		_ = version.Update().AddDownloads(1).Exec(c.Request().Context())
	}

	return c.Redirect(302, storage.GenerateDownloadLink(versionTarget.Key))
}

// @Summary Retrieve all Mod Versions
// @Tags Mod
// @Description Retrieve all mod versions by mod ID
// @Accept  json
// @Produce  json
// @Param modId path string true "Mod ID"
// @Success 200
// @Router /mod/{modId}/versions/all [get]
func getAllModVersions(c echo.Context) (interface{}, *ErrorResponse) {
	modID := c.Param("modId")

	mod, err := db.From(c.Request().Context()).Mod.Query().
		WithTags().
		Where(mod2.Or(mod2.ID(modID), mod2.ModReference(modID))).
		First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mod", slog.Any("err", err))
		return nil, &ErrorModNotFound
	}

	if mod == nil {
		return nil, &ErrorModNotFound
	}

	versions, err := mod.QueryVersions().
		WithVersionDependencies().
		WithTargets().
		Where(version2.Approved(true), version2.Denied(false)).
		Select(version2.FieldHash, version2.FieldSize, version2.FieldGameVersion, version2.FieldVersion).
		All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching versions", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	return (*conv.ModAllVersionsImpl)(nil).ConvertSlice(versions), nil
}
