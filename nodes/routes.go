package nodes

import (
	"github.com/labstack/echo/v4"
)

func RegisterModRoutes(router *echo.Group) {
	router.GET("/count", dataWrapper(getModCount))

	router.GET("/:modId", dataWrapper(getMod))

	router.GET("/:modId/latest-versions", dataWrapper(getModLatestVersions))
	router.GET("/:modId/versions", dataWrapper(getModVersions))
	router.GET("/:modId/authors", dataWrapper(getModAuthors))

	router.GET("/:modId/versions/all", dataWrapper(getAllModVersions))

	router.GET("/:modId/versions/:versionId", dataWrapper(getModVersion))
	router.GET("/:modId/versions/:versionId/download", downloadModVersion)
	router.HEAD("/:modId/versions/:versionId/download", downloadModVersion)
	router.GET("/:modId/versions/:versionId/:target/download", downloadModVersionTarget)
	router.HEAD("/:modId/versions/:versionId/:target/download", downloadModVersionTarget)
}

func RegisterModsRoutes(router *echo.Group) {
	router.GET("", dataWrapper(getMods))

	router.GET("/count", dataWrapper(getModCount))

	router.GET("/:modIds", dataWrapper(getModsByIDs))
	router.GET("/:modIds/latest-versions", dataWrapper(getModsLatestVersions))
}

func RegisterOAuthRoutes(router *echo.Group) {
	router.GET("/:url", dataWrapper(getOAuth))
	router.GET("/github", dataWrapper(getGithub))
}

func RegisterUserRoutes(router *echo.Group) {
	router.GET("/me", dataWrapper(authorized(getMe)))
	router.GET("/me/logout", dataWrapper(authorized(getLogout)))
	router.GET("/me/mods", dataWrapper(authorized(getMyMods)))

	router.GET("/:userId", dataWrapper(getUser))
	router.GET("/:userId/mods", dataWrapper(getUserMods))
}

func RegisterUsersRoutes(router *echo.Group) {
	router.GET("/:userIds", dataWrapper(getUsers))
}

func RegisterVersionRoutes(router *echo.Group) {
	router.GET("/:versionId", dataWrapper(getVersion))
	router.GET("/:versionId/download", downloadVersion)
	router.HEAD("/:versionId/download", downloadVersion)
	router.GET("/:versionId/:target/download", downloadModTarget)
	router.HEAD("/:versionId/:target/download", downloadModTarget)
}
