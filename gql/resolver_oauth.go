package gql

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/oauth"
	"github.com/satisfactorymodding/smr-api/util"
)

func (r *queryResolver) GetOAuthOptions(ctx context.Context, callbackURL string) (*generated.OAuthOptions, error) {
	unescapedURL, err := url.PathUnescape(callbackURL)
	if err != nil {
		return nil, fmt.Errorf("unable to unescape callback url: %w", err)
	}

	authOptions := oauth.GetOAuthOptions(unescapedURL)

	return &generated.OAuthOptions{
		Github:   authOptions["github"],
		Google:   authOptions["google"],
		Facebook: authOptions["facebook"],
	}, nil
}

func (r *mutationResolver) OAuthGithub(ctx context.Context, code string, state string) (*generated.UserSession, error) {
	if code == "" {
		return nil, errors.New("invalid oauth code")
	}

	u, err := oauth.GithubCallback(code, state)
	if err != nil {
		return nil, err
	}

	header := ctx.Value(util.ContextHeader{}).(http.Header)
	userAgent := header.Get("User-Agent")

	token, err := db.CompleteOAuthFlow(ctx, u, userAgent)
	if err != nil {
		return nil, err
	}

	return &generated.UserSession{
		Token: *token,
	}, nil
}

func (r *mutationResolver) OAuthGoogle(ctx context.Context, code string, state string) (*generated.UserSession, error) {
	if code == "" {
		return nil, errors.New("invalid oauth code")
	}

	u, err := oauth.GoogleCallback(code, state)
	if err != nil {
		return nil, err
	}

	header := ctx.Value(util.ContextHeader{}).(http.Header)
	userAgent := header.Get("User-Agent")

	token, err := db.CompleteOAuthFlow(ctx, u, userAgent)
	if err != nil {
		return nil, err
	}

	return &generated.UserSession{
		Token: *token,
	}, nil
}

func (r *mutationResolver) OAuthFacebook(ctx context.Context, code string, state string) (*generated.UserSession, error) {
	if code == "" {
		return nil, errors.New("invalid oauth code")
	}

	u, err := oauth.FacebookCallback(code, state)
	if err != nil {
		return nil, err
	}

	header := ctx.Value(util.ContextHeader{}).(http.Header)
	userAgent := header.Get("User-Agent")

	token, err := db.CompleteOAuthFlow(ctx, u, userAgent)
	if err != nil {
		return nil, err
	}

	return &generated.UserSession{
		Token: *token,
	}, nil
}
