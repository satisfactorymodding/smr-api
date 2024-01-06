package nodes

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db/postgres"
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

	version := postgres.GetVersion(c.Request().Context(), versionID)

	if version == nil {
		return nil, &ErrorVersionNotFound
	}

	return VersionToVersion(version), nil
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

	version := postgres.GetVersion(c.Request().Context(), versionID)

	if version == nil {
		return c.String(404, "version not found")
	}

	if redis.CanIncrement(c.RealIP(), "download", "version:"+versionID, time.Hour*4) {
		postgres.IncrementVersionDownloads(c.Request().Context(), version)
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

	version := postgres.GetVersion(c.Request().Context(), versionID)

	if version == nil {
		return c.String(404, "version not found, versionID:"+versionID)
	}

	versionTarget := postgres.GetVersionTarget(c.Request().Context(), versionID, target)

	if versionTarget == nil {
		return c.String(404, "target not found, versionID:"+versionID+" target:"+target)
	}

	if redis.CanIncrement(c.RealIP(), "download", "version:"+versionID, time.Hour*4) {
		postgres.IncrementVersionDownloads(c.Request().Context(), version)
	}

	return c.Redirect(302, storage.GenerateDownloadLink(versionTarget.Key))
}
