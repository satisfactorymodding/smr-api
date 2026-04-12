package nodes

import (
	"log/slog"
	"time"

	"github.com/Vilsol/slox"
	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/modpack"
	"github.com/satisfactorymodding/smr-api/generated/ent/modpackrelease"
	"github.com/satisfactorymodding/smr-api/models"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/util"
)

// @Summary Retrieve a list of Modpacks
// @Tags Modpacks
// @Description Retrieve a list of modpacks
// @Accept  json
// @Produce  json
// @Param limit query int false "How many modpacks to return"
// @Param offset query int false "Offset for list of modpacks to return"
// @Param order_by query string false "Order by field" Enums(created_at, updated_at, name, views, installs, hotness, popularity, search)
// @Param order query string false "Order of results" Enums(asc, desc)
// @Param search query string false "Search string"
// @Param hidden query bool false "Include hidden modpacks"
// @Success 200 {object} GenericResponse{data=[]generated.Modpack}
// @Router /modpacks [get]
func getModpacks(c echo.Context) (interface{}, *ErrorResponse) {
	limit := util.GetIntRange(c, "limit", 1, 100, 25)
	offset := util.GetIntRange(c, "offset", 0, 9999999, 0)
	orderBy := util.OneOf(c, "order_by", []string{"created_at", "updated_at", "name", "views", "installs", "hotness", "popularity", "search"}, "created_at")
	order := util.OneOf(c, "order", []string{"asc", "desc"}, "desc")
	search := c.QueryParam("search")
	hiddenParam := c.QueryParam("hidden")

	modpackFilter := models.DefaultModpackFilter()
	modpackFilter.Limit = &limit
	modpackFilter.Offset = &offset

	orderByGen := generated.ModpackFields(orderBy)
	modpackFilter.OrderBy = &orderByGen

	orderGen := generated.Order(order)
	modpackFilter.Order = &orderGen

	if search != "" {
		modpackFilter.Search = &search
	}

	if hiddenParam != "" {
		hidden := hiddenParam == "true"
		modpackFilter.Hidden = &hidden
	}

	query := db.From(c.Request().Context()).Modpack.Query().
		WithTags().
		WithModpackMods()

	query = db.ConvertModpackFilter(query, modpackFilter, false)

	result, err := query.All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed to get modpacks", slog.Any("error", err))
		return nil, GenericUserError(err)
	}

	return (*conv.ModpackImpl)(nil).ConvertSlice(result), nil
}

// @Summary Retrieve a single Modpack
// @Tags Modpacks
// @Description Retrieve a single modpack by ID
// @Accept  json
// @Produce  json
// @Param modpackId path string true "Modpack ID"
// @Success 200 {object} GenericResponse{data=generated.Modpack}
// @Failure 404 {object} GenericResponse{error=ErrorResponse}
// @Router /modpacks/{modpackId} [get]
func getModpack(c echo.Context) (interface{}, *ErrorResponse) {
	modpackID := c.Param("modpackId")

	dbModpack, err := db.From(c.Request().Context()).Modpack.Query().
		Where(modpack.ID(modpackID)).
		WithTags().
		WithReleases(func(q *ent.ModpackReleaseQuery) {
			q.WithTargets()
		}).
		WithModpackMods().
		WithParent().
		WithChildren().
		First(c.Request().Context())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ErrorModpackNotFound
		}
		slox.Error(c.Request().Context(), "failed to get modpack", slog.Any("error", err))
		return nil, GenericUserError(err)
	}

	if dbModpack == nil {
		return nil, &ErrorModpackNotFound
	}

	// Increment view count with rate limiting
	if redis.CanIncrement(c.RealIP(), "view", "modpack:"+modpackID, time.Hour*4) {
		if err := dbModpack.Update().AddViews(1).Exec(c.Request().Context()); err != nil {
			slox.Error(c.Request().Context(), "failed to increment modpack views", slog.Any("error", err))
		}
	}

	return (*conv.ModpackImpl)(nil).Convert(dbModpack), nil
}

// @Summary Retrieve a Modpack Release
// @Tags Modpacks
// @Description Retrieve a specific release of a modpack
// @Accept  json
// @Produce  json
// @Param modpackId path string true "Modpack ID"
// @Param version path string true "Release Version"
// @Success 200 {object} GenericResponse{data=generated.ModpackRelease}
// @Failure 404 {object} GenericResponse{error=ErrorResponse}
// @Router /modpacks/{modpackId}/releases/{version} [get]
func getModpackRelease(c echo.Context) (interface{}, *ErrorResponse) {
	modpackID := c.Param("modpackId")
	version := c.Param("version")

	dbRelease, err := db.From(c.Request().Context()).ModpackRelease.Query().
		Where(modpackrelease.HasModpackWith(modpack.ID(modpackID)), modpackrelease.Version(version)).
		WithTargets().
		First(c.Request().Context())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ErrorModpackReleaseNotFound
		}
		slox.Error(c.Request().Context(), "failed to get modpack release", slog.Any("error", err))
		return nil, GenericUserError(err)
	}

	return (*conv.ModpackReleaseImpl)(nil).Convert(dbRelease), nil
}
