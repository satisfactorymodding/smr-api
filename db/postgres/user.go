package postgres

import (
	"context"

	"github.com/satisfactorymodding/smr-api/oauth"
	"github.com/satisfactorymodding/smr-api/util"
)

func GetUserSession(ctx context.Context, oauthUser *oauth.UserData, userAgent string) (*UserSession, *User, bool) {
	user := User{
		Email:      oauthUser.Email,
		Avatar:     oauthUser.Avatar,
		JoinedFrom: string(oauthUser.Site),
		Username:   oauthUser.Username,
	}

	// Find or create the user by email
	find := DBCtx(ctx).Where(&User{Email: oauthUser.Email})

	if oauthUser.Site == oauth.SiteGithub {
		find = find.Or(&User{GithubID: &oauthUser.ID})
	} else if oauthUser.Site == oauth.SiteGoogle {
		find = find.Or(&User{GoogleID: &oauthUser.ID})
	} else if oauthUser.Site == oauth.SiteFacebook {
		find = find.Or(&User{FacebookID: &oauthUser.ID})
	}

	find.Find(&user)

	newUser := false

	if user.ID == "" {
		user.ID = util.GenerateUniqueID()

		if oauthUser.Site == oauth.SiteGithub {
			user.GithubID = &oauthUser.ID
		} else if oauthUser.Site == oauth.SiteGoogle {
			user.GoogleID = &oauthUser.ID
		} else if oauthUser.Site == oauth.SiteFacebook {
			user.FacebookID = &oauthUser.ID
		}

		DBCtx(ctx).Create(&user)
		newUser = true
	}

	if !newUser {
		newID := false
		if oauthUser.Site == oauth.SiteGithub && user.GithubID == nil {
			user.GithubID = &oauthUser.ID
			newID = true
		} else if oauthUser.Site == oauth.SiteGoogle && user.GoogleID == nil {
			user.GoogleID = &oauthUser.ID
			newID = true
		} else if oauthUser.Site == oauth.SiteFacebook && user.FacebookID == nil {
			user.FacebookID = &oauthUser.ID
			newID = true
		}

		if newID {
			Save(ctx, &user)
		}
	}

	// TODO Archive old deleted sessions to cold storage
	// DBCtx(ctx).Delete(&UserSession{UserAgent: userAgent})

	session := UserSession{
		User:      user,
		Token:     util.GenerateUserToken(),
		UserAgent: userAgent,
	}

	session.ID = util.GenerateUniqueID()

	// Create a new session
	DBCtx(ctx).Create(&session)

	return &session, &user, newUser
}

func LogoutSession(ctx context.Context, token string) {
	// TODO Archive old deleted sessions to cold storage
	DBCtx(ctx).Delete(&UserSession{Token: token})
}

func GetUserByToken(ctx context.Context, token string) *User {
	// TODO Merge into a single query
	var session UserSession
	DBCtx(ctx).Find(&session, UserSession{Token: token})

	if session.ID == "" {
		return nil
	}

	var user User
	DBCtx(ctx).Find(&user, "id = ?", session.UserID)

	if user.ID == "" {
		return nil
	}

	return &user
}

func GetUserByID(ctx context.Context, userID string) *User {
	var user User
	DBCtx(ctx).Find(&user, "id = ?", userID)

	if user.ID == "" {
		return nil
	}

	return &user
}

func GetUsersByID(ctx context.Context, userIds []string) *[]User {
	var users []User
	DBCtx(ctx).Find(&users, "id in (?)", userIds)

	if len(userIds) != len(users) {
		return nil
	}

	return &users
}

func GetUserMods(ctx context.Context, userID string) []UserMod {
	var mods []UserMod
	DBCtx(ctx).Raw("SELECT * from \"user_mods\" as tdm WHERE user_id = ? AND mod_id = (SELECT id FROM mods WHERE id = tdm.mod_id AND deleted_at is NULL LIMIT 1)", userID).Find(&mods)
	return mods
}

func GetModAuthors(ctx context.Context, modID string) []UserMod {
	var authors []UserMod
	DBCtx(ctx).Find(&authors, "mod_id = ?", modID)
	return authors
}
