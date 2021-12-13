package gql

import (
	"context"
	"net/http"
	"reflect"

	"github.com/satisfactorymodding/smr-api/auth"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/util"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
)

func MakeDirective() generated.DirectiveRoot {
	return generated.DirectiveRoot{
		CanEditGuide:             canEditGuide,
		CanEditMod:               canEditMod,
		CanEditVersion:           canEditVersion,
		IsLoggedIn:               isLoggedIn,
		IsNotLoggedIn:            isNotLoggedIn,
		CanEditUser:              canEditUser,
		CanApproveMods:           canApproveMods,
		CanApproveVersions:       canApproveVersions,
		CanEditUsers:             canEditUsers,
		CanEditSMLVersions:       canEditSMLVersions,
		CanEditBootstrapVersions: canEditBootstrapVersions,
		CanEditAnnouncements:     canEditAnnouncements,
	}
}

type Directive struct {
	generated.DirectiveRoot
}

func canEditMod(ctx context.Context, obj interface{}, next graphql.Resolver, field string) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	dbMod := postgres.GetModByID(ctx, getArgument(ctx, field).(string))

	if dbMod == nil {
		return nil, errors.New("mod not found")
	}

	if postgres.UserCanUploadModVersions(ctx, user, dbMod.ID) {
		return next(ctx)
	}

	if user.Has(ctx, auth.RoleEditAnyContent) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditVersion(ctx context.Context, obj interface{}, next graphql.Resolver, field string) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	dbVersion := postgres.GetVersion(ctx, getArgument(ctx, field).(string))

	if dbVersion == nil {
		return nil, errors.New("version not found")
	}

	if postgres.UserCanUploadModVersions(ctx, user, dbVersion.ModID) {
		return next(ctx)
	}

	if user.Has(ctx, auth.RoleEditAnyContent) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditUser(ctx context.Context, obj interface{}, next graphql.Resolver, field string, object bool) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	var userID string
	if object {
		userID = reflect.Indirect(reflect.ValueOf(obj)).FieldByName(field).String()
	} else {
		userID = getArgument(ctx, field).(string)
	}

	dbUser := postgres.GetUserByID(ctx, userID)

	if dbUser == nil {
		return nil, errors.New("user not found")
	}

	if dbUser.ID == user.ID {
		return next(ctx)
	}

	if user.Has(ctx, auth.RoleEditUsers) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditGuide(ctx context.Context, obj interface{}, next graphql.Resolver, field string) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	dbGuide := postgres.GetGuideByID(ctx, getArgument(ctx, field).(string))

	if dbGuide == nil {
		return nil, errors.New("guide not found")
	}

	if dbGuide.UserID == user.ID {
		return next(ctx)
	}

	if user.Has(ctx, auth.RoleEditAnyContent) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

var errUserNotLoggedIn = errors.New("user not logged in")

func isLoggedIn(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	header := ctx.Value(util.ContextHeader{}).(http.Header)
	authorization := header.Get("Authorization")

	if authorization == "" {
		return nil, errUserNotLoggedIn
	}

	payload, err := util.VerifyUserToken(authorization)
	if err != nil {
		return nil, errors.New("invalid user authorization token")
	}

	userID := payload.Get("userID")
	if userID == "" {
		return nil, errUserNotLoggedIn
	}

	if redis.IsAccessTokenRevoked(authorization) {
		return nil, errUserNotLoggedIn
	}

	user := postgres.GetUserByID(ctx, userID)
	if user == nil {
		return nil, errUserNotLoggedIn
	}

	if user.Banned {
		return nil, errors.New("user banned")
	}

	userCtx := context.WithValue(ctx, postgres.UserKey{}, user)

	return next(userCtx)
}

func isNotLoggedIn(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	header := ctx.Value(util.ContextHeader{}).(http.Header)
	authorization := header.Get("Authorization")

	if authorization != "" {
		payload, _ := util.VerifyUserToken(authorization)
		userID := payload.Get("userID")

		if userID != "" {
			return nil, errors.New("user is logged in")
		}
	}

	return next(ctx)
}

func getArgument(ctx context.Context, key string) interface{} {
	return graphql.GetFieldContext(ctx).Args[key]
}

func canApproveMods(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleApproveMods) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canApproveVersions(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleApproveVersions) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditUsers(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleEditUsers) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditSMLVersions(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleEditSMLVersions) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditBootstrapVersions(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleEditBootstrapVersions) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditAnnouncements(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleEditAnnouncements) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}
