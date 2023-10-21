package gql

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/oauth"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
)

func (r *queryResolver) GetOAuthOptions(ctx context.Context, callbackURL string) (*generated.OAuthOptions, error) {
	wrapper, _ := WrapQueryTrace(ctx, "getOAuthOptions")
	defer wrapper.end()

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
	wrapper, newCtx := WrapMutationTrace(ctx, "oAuthGithub")
	defer wrapper.end()

	if code == "" {
		return nil, errors.New("invalid oauth code")
	}

	user, err := oauth.GithubCallback(code, state)
	if err != nil {
		return nil, err
	}

	header := ctx.Value(util.ContextHeader{}).(http.Header)
	userAgent := header.Get("User-Agent")

	token, err := completeOAuthFlow(newCtx, user, userAgent)
	if err != nil {
		return nil, err
	}

	return &generated.UserSession{
		Token: *token,
	}, nil
}

func (r *mutationResolver) OAuthGoogle(ctx context.Context, code string, state string) (*generated.UserSession, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "oAuthGoogle")
	defer wrapper.end()

	if code == "" {
		return nil, errors.New("invalid oauth code")
	}

	user, err := oauth.GoogleCallback(code, state)
	if err != nil {
		return nil, err
	}

	header := ctx.Value(util.ContextHeader{}).(http.Header)
	userAgent := header.Get("User-Agent")

	token, err := completeOAuthFlow(newCtx, user, userAgent)
	if err != nil {
		return nil, err
	}

	return &generated.UserSession{
		Token: *token,
	}, nil
}

func (r *mutationResolver) OAuthFacebook(ctx context.Context, code string, state string) (*generated.UserSession, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "oAuthFacebook")
	defer wrapper.end()

	if code == "" {
		return nil, errors.New("invalid oauth code")
	}

	user, err := oauth.FacebookCallback(code, state)
	if err != nil {
		return nil, err
	}

	header := ctx.Value(util.ContextHeader{}).(http.Header)
	userAgent := header.Get("User-Agent")

	token, err := completeOAuthFlow(newCtx, user, userAgent)
	if err != nil {
		return nil, err
	}

	return &generated.UserSession{
		Token: *token,
	}, nil
}

func completeOAuthFlow(ctx context.Context, user *oauth.UserData, userAgent string) (*string, error) {
	avatarURL := user.Avatar
	user.Avatar = ""

	session, dbUser, newUser := postgres.GetUserSession(ctx, user, userAgent)

	if avatarURL != "" && newUser {
		avatarData, err := util.LinkToWebp(ctx, avatarURL)
		if err != nil {
			return nil, err
		}

		success, avatarKey := storage.UploadUserAvatar(ctx, session.UserID, bytes.NewReader(avatarData))
		if success {
			dbUser.Avatar = storage.GenerateDownloadLink(avatarKey)
			postgres.Save(ctx, &dbUser)
		}
	}

	return &session.Token, nil
}
