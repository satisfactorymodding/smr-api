package nodes

import (
	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/generated/ent"
)

type DataFunction func(c echo.Context) (data interface{}, err *ErrorResponse)

func dataWrapper(nested DataFunction) func(c echo.Context) error {
	return func(c echo.Context) error {
		data, err := nested(c)
		if err != nil {
			return c.JSON(err.Status, GenericResponse{
				Success: false,
				Error:   err,
			})
		}

		return c.JSON(200, GenericResponse{
			Success: true,
			Data:    data,
		})
	}
}

type AuthorizedDataFunction func(user *ent.User, c echo.Context) (data interface{}, err *ErrorResponse)

func authorized(nested AuthorizedDataFunction) DataFunction {
	return func(c echo.Context) (interface{}, *ErrorResponse) {
		user := userFromContext(c)

		if user == nil {
			return nil, &ErrorInvalidAuthorizationToken
		}

		if user.Banned {
			return nil, &ErrorUserBanned
		}

		return nested(user, c)
	}
}
