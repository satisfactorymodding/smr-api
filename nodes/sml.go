package nodes

import (
	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db/postgres"
)

// @Summary Retrieve a list of latest versions for sml
// @Tags SML
// @Description Retrieve a list of latest versions for sml
// @Accept  json
// @Produce  json
// @Success 200
// @Router /sml/latest-versions [get]
func getSMLLatestVersions(c echo.Context) (interface{}, *ErrorResponse) {
	smlVersions := postgres.GetSMLLatestVersions(c.Request().Context())

	if smlVersions == nil {
		return nil, &ErrorVersionNotFound
	}

	result := make(map[string]*SMLVersion)

	for _, v := range *smlVersions {
		result[v.Stability] = SMLVersionToSMLVersion(&v)
	}

	return result, nil
}
