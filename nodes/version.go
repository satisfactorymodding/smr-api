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

// @Summary Download a Platform
// @Tags Version
// @Tags Platform
// @Description Download a mod version by version ID and Platform
// @Accept  json
// @Produce  json
// @Param versionId path string true "Version ID"
// @Param versionId path string true "Version ID"
// @Success 200
// @Router /versions/{versionId}/{platform}/download [get]
func downloadModArch(c echo.Context) error {
	versionID := c.Param("versionId")
	platform := c.Param("platform")

	version := postgres.GetVersion(c.Request().Context(), versionID)

	if version == nil {
		return c.String(404, "version not found, versionID:"+versionID)
	}

	arch := postgres.GetModArchByPlatform(c.Request().Context(), versionID, platform)

	if arch == nil {
		return c.String(404, "platform not found, versionID:"+versionID+" platform:"+platform)
	}

	if redis.CanIncrement(c.RealIP(), "download", "version:"+versionID, time.Hour*4) {
		postgres.IncrementVersionDownloads(c.Request().Context(), version)
	}

	return c.Redirect(302, storage.GenerateDownloadLink(arch.Key))
}
