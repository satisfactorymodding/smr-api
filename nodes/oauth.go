package nodes

import (
	"bytes"
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/oauth"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
)

// @Summary Retrieve a list of OAuth methods
// @Tags OAuth
// @Accept  json
// @Produce  json
// @Success 200
// @Router /oauth [get]
func getOAuth(c echo.Context) (interface{}, *ErrorResponse) {
	callbackURL := c.Param("url")
	unescapedURL, err := url.PathUnescape(callbackURL)
	if err != nil {
		return nil, GenericUserError(err)
	}

	return oauth.GetOAuthOptions(unescapedURL), nil
}

// @Summary Callback URL for github OAuth
// @Tags OAuth
// @Accept  json
// @Produce  json
// @Param code query string true "OAuth Code"
// @Param state query string true "OAuth Code"
// @Success 200
// @Router /oauth/github [get]
func getGithub(c echo.Context) (interface{}, *ErrorResponse) {
	code := c.QueryParam("code")

	if code == "" {
		return nil, &ErrorInvalidOAuthCode
	}

	state := c.QueryParam("state")

	if state == "" {
		return nil, &ErrorInvalidOAuthCode
	}

	user, err := oauth.GithubCallback(code, state)
	if err != nil {
		return nil, GenericUserError(err)
	}

	userAgent := c.Request().Header.Get("User-Agent")

	avatarURL := user.Avatar
	user.Avatar = ""

	session, dbUser, newUser := postgres.GetUserSession(c.Request().Context(), user, userAgent)

	if avatarURL != "" && newUser {
		avatarData, err := util.LinkToWebp(c.Request().Context(), avatarURL)
		if err != nil {
			return nil, GenericUserError(err)
		}

		success, avatarKey := storage.UploadUserAvatar(c.Request().Context(), session.UserID, bytes.NewReader(avatarData))
		if success {
			dbUser.Avatar = storage.GenerateDownloadLink(avatarKey)
			postgres.Save(c.Request().Context(), &dbUser)
		}
	}

	return SessionToSession(session), nil
}
