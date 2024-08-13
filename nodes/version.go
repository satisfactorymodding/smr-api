package nodes

import (
	"log/slog"
	"time"

	"github.com/Vilsol/slox"
	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent/versiontarget"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/storage"
)

// @Summary Retrieve a Version
// @Tags Version
// @Description Retrieve a version by version ID
// @Accept  json
// @Produce  json
// @Param versionId path string true "Version ID"
// @Success 200
// @Router /version/{versionId} [get]
func getVersion(c echo.Context) (interface{}, *ErrorResponse) {
	versionID := c.Param("versionId")

	version, err := db.From(c.Request().Context()).Version.Get(c.Request().Context(), versionID)
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching version", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	if version == nil {
		return nil, &ErrorVersionNotFound
	}

	return (*conv.VersionImpl)(nil).Convert(version), nil
}

// @Summary Download a Version
// @Tags Version
// @Description Download a mod version by version ID
// @Accept  json
// @Produce  json
// @Param versionId path string true "Version ID"
// @Success 200
// @Router /versions/{versionId}/download [get]
func downloadVersion(c echo.Context) error {
	versionID := c.Param("versionId")

	version, err := db.From(c.Request().Context()).Version.Get(c.Request().Context(), versionID)
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching version", slog.Any("err", err))
		return err
	}

	if version == nil {
		return c.String(404, "version not found")
	}

	if c.Request().Method == echo.GET &&
		redis.CanIncrement(c.RealIP(), "download", "version:"+versionID, time.Hour*4) {
		_ = version.Update().AddDownloads(1).Exec(c.Request().Context())
	}

	return c.Redirect(302, storage.GenerateDownloadLink(version.Key))
}

// @Summary Download a TargetName
// @Tags Version
// @Tags TargetName
// @Description Download a mod version by version ID and TargetName
// @Accept  json
// @Produce  json
// @Param versionId path string true "Version ID"
// @Param target path string true "TargetName"
// @Success 200
// @Router /versions/{versionId}/{target}/download [get]
func downloadModTarget(c echo.Context) error {
	versionID := c.Param("versionId")
	target := c.Param("target")

	version, err := db.From(c.Request().Context()).Version.Get(c.Request().Context(), versionID)
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching version", slog.Any("err", err))
		return err
	}

	if version == nil {
		return c.String(404, "version not found, versionID:"+versionID)
	}

	versionTarget, err := version.QueryTargets().Where(versiontarget.TargetName(target)).First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching target", slog.Any("err", err))
		return err
	}

	if versionTarget == nil {
		return c.String(404, "target not found, versionID:"+versionID+" target:"+target)
	}

	if c.Request().Method == echo.GET &&
		redis.CanIncrement(c.RealIP(), "download", "version:"+versionID, time.Hour*4) {
		_ = version.Update().AddDownloads(1).Exec(c.Request().Context())
	}

	return c.Redirect(302, storage.GenerateDownloadLink(versionTarget.Key))
}
