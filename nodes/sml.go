package nodes

import (
	"log/slog"

	"entgo.io/ent/dialect/sql"
	"github.com/Vilsol/slox"
	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent/smlversion"
)

// @Summary Retrieve a list of latest versions for sml
// @Tags SML
// @Description Retrieve a list of latest versions for sml
// @Accept  json
// @Produce  json
// @Success 200
// @Router /sml/latest-versions [get]
func getSMLLatestVersions(c echo.Context) (interface{}, *ErrorResponse) {
	smlVersions, err := db.From(c.Request().Context()).SmlVersion.Query().
		WithTargets().
		Modify(func(s *sql.Selector) {
			s.SelectExpr(sql.ExprP("distinct on (stability) *"))
		}).
		Order(smlversion.ByStability(sql.OrderDesc()), smlversion.ByCreatedAt(sql.OrderDesc())).
		All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching sml versions", slog.Any("err", err))
		return nil, &ErrorVersionNotFound
	}

	return (*conv.SMLVersionImpl)(nil).ConvertSlice(smlVersions), nil
}
