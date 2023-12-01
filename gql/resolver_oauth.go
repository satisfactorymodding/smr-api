package gql

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
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
	wrapper, ctx := WrapMutationTrace(ctx, "oAuthGithub")
	defer wrapper.end()

	if code == "" {
		return nil, errors.New("invalid oauth code")
	}

	u, err := oauth.GithubCallback(code, state)
	if err != nil {
		return nil, err
	}

	header := ctx.Value(util.ContextHeader{}).(http.Header)
	userAgent := header.Get("User-Agent")

	token, err := completeOAuthFlow(ctx, u, userAgent)
	if err != nil {
		return nil, err
	}

	return &generated.UserSession{
		Token: *token,
	}, nil
}

func (r *mutationResolver) OAuthGoogle(ctx context.Context, code string, state string) (*generated.UserSession, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "oAuthGoogle")
	defer wrapper.end()

	if code == "" {
		return nil, errors.New("invalid oauth code")
	}

	u, err := oauth.GoogleCallback(code, state)
	if err != nil {
		return nil, err
	}

	header := ctx.Value(util.ContextHeader{}).(http.Header)
	userAgent := header.Get("User-Agent")

	token, err := completeOAuthFlow(ctx, u, userAgent)
	if err != nil {
		return nil, err
	}

	return &generated.UserSession{
		Token: *token,
	}, nil
}

func (r *mutationResolver) OAuthFacebook(ctx context.Context, code string, state string) (*generated.UserSession, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "oAuthFacebook")
	defer wrapper.end()

	if code == "" {
		return nil, errors.New("invalid oauth code")
	}

	u, err := oauth.FacebookCallback(code, state)
	if err != nil {
		return nil, err
	}

	header := ctx.Value(util.ContextHeader{}).(http.Header)
	userAgent := header.Get("User-Agent")

	token, err := completeOAuthFlow(ctx, u, userAgent)
	if err != nil {
		return nil, err
	}

	return &generated.UserSession{
		Token: *token,
	}, nil
}

func completeOAuthFlow(ctx context.Context, u *oauth.UserData, userAgent string) (*string, error) {
	avatarURL := u.Avatar
	u.Avatar = ""

	find := db.From(ctx).User.Query().Where(user.Email(u.Email))

	if u.Site == oauth.SiteGithub {
		find = find.Where(user.GithubID(u.ID))
	} else if u.Site == oauth.SiteGoogle {
		find = find.Where(user.GoogleID(u.ID))
	} else if u.Site == oauth.SiteFacebook {
		find = find.Where(user.FacebookID(u.ID))
	}

	found, err := find.First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	newUser := false
	if ent.IsNotFound(err) {
		var err error
		create := db.From(ctx).User.
			Create().
			SetEmail(u.Email).
			SetAvatar(u.Avatar).
			SetJoinedFrom(string(u.Site)).
			SetUsername(u.Username)

		if u.Site == oauth.SiteGithub {
			create = create.SetGithubID(u.ID)
		} else if u.Site == oauth.SiteGoogle {
			create = create.SetGoogleID(u.ID)
		} else if u.Site == oauth.SiteFacebook {
			create = create.SetFacebookID(u.ID)
		}

		found, err = create.Save(ctx)
		if err != nil {
			return nil, err
		}

		newUser = true
	}

	if !newUser {
		var update *ent.UserUpdateOne
		if u.Site == oauth.SiteGithub && found.GithubID == "" {
			update = found.Update().SetGithubID(u.ID)
		} else if u.Site == oauth.SiteGoogle && found.GoogleID == "" {
			update = found.Update().SetGoogleID(u.ID)
		} else if u.Site == oauth.SiteFacebook && found.FacebookID == "" {
			update = found.Update().SetFacebookID(u.ID)
		}

		if update != nil {
			if err := update.Exec(ctx); err != nil {
				return nil, err
			}
		}
	}

	// TODO Archive old deleted sessions to cold storage

	session, err := db.From(ctx).UserSession.
		Create().
		SetUserID(found.ID).
		SetToken(util.GenerateUserToken()).
		SetUserAgent(userAgent).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if avatarURL != "" && newUser {
		avatarData, err := util.LinkToWebp(ctx, avatarURL)
		if err != nil {
			return nil, err
		}

		success, avatarKey := storage.UploadUserAvatar(ctx, found.ID, bytes.NewReader(avatarData))
		if success {
			if err := found.Update().SetAvatar(storage.GenerateDownloadLink(avatarKey)).Exec(ctx); err != nil {
				return nil, err
			}
		}
	}

	return &session.Token, nil
}
