package nodes

import (
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/oauth"
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

	token, err := db.CompleteOAuthFlow(c.Request().Context(), user, userAgent)
	if err != nil {
		return nil, GenericUserError(err)
	}

	return &UserSession{
		Token: *token,
	}, nil
}
