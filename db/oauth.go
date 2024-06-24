package db

import (
	"bytes"
	"context"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"

	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/oauth"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
)

func CompleteOAuthFlow(ctx context.Context, u *oauth.UserData, userAgent string) (*string, error) {
	avatarURL := u.Avatar
	u.Avatar = ""

	var oauthPredicate predicate.User

	if u.Site == oauth.SiteGithub {
		oauthPredicate = user.GithubID(u.ID)
	} else if u.Site == oauth.SiteGoogle {
		oauthPredicate = user.GoogleID(u.ID)
	} else if u.Site == oauth.SiteFacebook {
		oauthPredicate = user.FacebookID(u.ID)
	}

	find := From(ctx).User.Query().Where(user.Or(user.Email(u.Email), oauthPredicate))

	found, err := find.First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	newUser := false
	if ent.IsNotFound(err) {
		create := From(ctx).User.
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
			if err = update.Exec(ctx); err != nil {
				return nil, err
			}
		}
	}

	// TODO Archive old deleted sessions to cold storage

	session, err := From(ctx).UserSession.
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
