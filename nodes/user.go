package nodes

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db/postgres"
)

func userFromContext(c echo.Context) *postgres.User {
	authorization := c.Request().Header.Get("Authorization")

	if authorization == "" {
		return nil
	}

	user := postgres.GetUserByToken(c.Request().Context(), authorization)

	if user == nil {
		return nil
	}

	return user
}

// @Summary Retrieve Current User
// @Tags User
// @Description Retrieve the user associated with the token
// @Accept  json
// @Produce  json
// @Success 200
// @Router /user/me [get]
func getMe(user *postgres.User, _ echo.Context) (interface{}, *ErrorResponse) {
	return UserToPrivateUser(user), nil
}

// @Summary Log Out Current User
// @Tags User
// @Description Log out the user associated with the token
// @Accept  json
// @Produce  json
// @Success 200
// @Router /user/me/logout [get]
func getLogout(_ *postgres.User, c echo.Context) (interface{}, *ErrorResponse) {
	postgres.LogoutSession(c.Request().Context(), c.Request().Header.Get("Authorization"))
	return nil, nil
}

// @Summary Retrieve Current Users Mods
// @Tags User
// @Description Retrieve the users mods associated with the token
// @Accept  json
// @Produce  json
// @Success 200
// @Router /user/me/mods [get]
func getMyMods(user *postgres.User, c echo.Context) (interface{}, *ErrorResponse) {
	mods := postgres.GetUserMods(c.Request().Context(), user.ID)

	converted := make([]*UserMod, len(mods))
	for k, v := range mods {
		converted[k] = UserModToUserMod(&v)
	}

	return converted, nil
}

// @Summary Retrieve a list of Users
// @Tags Users
// @Description Retrieve a list of users by user ID
// @Accept  json
// @Produce  json
// @Success 200
// @Param userIds path string true "User IDs comma-separated"
// @Success 200
// @Router /users/{userIds} [get]
func getUsers(c echo.Context) (interface{}, *ErrorResponse) {
	userID := c.Param("userIds")
	userIDSplit := strings.Split(userID, ",")

	// TODO limit amount of users requestable

	users := postgres.GetUsersByID(c.Request().Context(), userIDSplit)

	if users == nil {
		return nil, &ErrorUserNotFound
	}

	converted := make([]*PublicUser, len(*users))
	for k, v := range *users {
		converted[k] = UserToPublicUser(&v)
	}

	return converted, nil
}

// @Summary Retrieve a Users Mods
// @Tags User
// @Description Retrieve a users mods by user ID
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Success 200
// @Router /user/{userId}/mods [get]
func getUserMods(c echo.Context) (interface{}, *ErrorResponse) {
	userID := c.Param("userId")

	user := postgres.GetUserByID(c.Request().Context(), userID)

	if user == nil {
		return nil, &ErrorUserNotFound
	}

	mods := postgres.GetUserMods(c.Request().Context(), user.ID)

	converted := make([]*UserMod, len(mods))
	for k, v := range mods {
		converted[k] = UserModToUserMod(&v)
	}

	return converted, nil
}

// @Summary Retrieve a User
// @Tags User
// @Description Retrieve a user by user ID
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Success 200
// @Router /user/{userId} [get]
func getUser(c echo.Context) (interface{}, *ErrorResponse) {
	userID := c.Param("userId")

	user := postgres.GetUserByID(c.Request().Context(), userID)

	if user == nil {
		return nil, &ErrorUserNotFound
	}

	return UserToPublicUser(user), nil
}
