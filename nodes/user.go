package nodes

import (
	"log/slog"
	"strings"

	"github.com/Vilsol/slox"
	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/generated/ent/usermod"
	"github.com/satisfactorymodding/smr-api/generated/ent/usersession"
)

func userFromContext(c echo.Context) *ent.User {
	authorization := c.Request().Header.Get("Authorization")

	if authorization == "" {
		return nil
	}

	user, err := db.From(c.Request().Context()).User.Query().
		Where(user.HasSessionsWith(usersession.Token(authorization))).
		First(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mods", slog.Any("err", err))
		return nil
	}

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
func getMe(user *ent.User, _ echo.Context) (interface{}, *ErrorResponse) {
	return &User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
	}, nil
}

// @Summary Log Out Current User
// @Tags User
// @Description Log out the user associated with the token
// @Accept  json
// @Produce  json
// @Success 200
// @Router /user/me/logout [get]
func getLogout(_ *ent.User, c echo.Context) (interface{}, *ErrorResponse) {
	if _, err := db.From(c.Request().Context()).UserSession.Delete().
		Where(usersession.Token(c.Request().Header.Get("Authorization"))).
		Exec(c.Request().Context()); err != nil {
		slox.Error(c.Request().Context(), "failed deleting session", slog.Any("err", err))
		return nil, &ErrorUserNotFound
	}
	return nil, nil
}

// @Summary Retrieve Current Users Mods
// @Tags User
// @Description Retrieve the users mods associated with the token
// @Accept  json
// @Produce  json
// @Success 200
// @Router /user/me/mods [get]
func getMyMods(user *ent.User, c echo.Context) (interface{}, *ErrorResponse) {
	mods, err := db.From(c.Request().Context()).UserMod.Query().Where(usermod.UserID(user.ID)).All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching authors", slog.Any("err", err))
		return nil, &ErrorUserNotFound
	}

	return (*conv.UserModImpl)(nil).ConvertSlice(mods), nil
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

	users, err := db.From(c.Request().Context()).User.Query().Where(user.IDIn(userIDSplit...)).All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching users", slog.Any("err", err))
		return nil, nil
	}

	return (*conv.UserImpl)(nil).ConvertSlice(users), nil
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

	user, err := db.From(c.Request().Context()).User.Get(c.Request().Context(), userID)
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mods", slog.Any("err", err))
		return nil, &ErrorUserNotFound
	}

	if user == nil {
		return nil, &ErrorUserNotFound
	}

	mods, err := db.From(c.Request().Context()).UserMod.Query().Where(usermod.UserID(user.ID)).All(c.Request().Context())
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mods", slog.Any("err", err))
		return nil, &ErrorModNotFound
	}

	return (*conv.UserModImpl)(nil).ConvertSlice(mods), nil
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

	user, err := db.From(c.Request().Context()).User.Get(c.Request().Context(), userID)
	if err != nil {
		slox.Error(c.Request().Context(), "failed fetching mods", slog.Any("err", err))
		return nil, &ErrorUserNotFound
	}

	if user == nil {
		return nil, &ErrorUserNotFound
	}

	return &PublicUser{
		ID:        user.ID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
	}, nil
}
