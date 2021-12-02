package gql

import (
	"context"
	"net/http"
	"reflect"

	"github.com/satisfactorymodding/smr-api/auth"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
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

	dbMod := postgres.GetModByID(getArgument(ctx, field).(string), &ctx)

	if dbMod == nil {
		return nil, errors.New("mod not found")
	}

	if postgres.UserCanUploadModVersions(user, dbMod.ID, &ctx) {
		return next(ctx)
	}

	if user.Has(auth.RoleEditAnyContent, &ctx) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditVersion(ctx context.Context, obj interface{}, next graphql.Resolver, field string) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	dbVersion := postgres.GetVersion(getArgument(ctx, field).(string), &ctx)

	if dbVersion == nil {
		return nil, errors.New("version not found")
	}

	if postgres.UserCanUploadModVersions(user, dbVersion.ModID, &ctx) {
		return next(ctx)
	}

	if user.Has(auth.RoleEditAnyContent, &ctx) {
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

	dbUser := postgres.GetUserByID(userID, &ctx)

	if dbUser == nil {
		return nil, errors.New("user not found")
	}

	if dbUser.ID == user.ID {
		return next(ctx)
	}

	if user.Has(auth.RoleEditUsers, &ctx) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditGuide(ctx context.Context, obj interface{}, next graphql.Resolver, field string) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	dbGuide := postgres.GetGuideByID(getArgument(ctx, field).(string), &ctx)

	if dbGuide == nil {
		return nil, errors.New("guide not found")
	}

	if dbGuide.UserID == user.ID {
		return next(ctx)
	}

	if user.Has(auth.RoleEditAnyContent, &ctx) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func isLoggedIn(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	header := ctx.Value(util.ContextHeader{}).(http.Header)
	authorization := header.Get("Authorization")

	if authorization == "" {
		return nil, errors.New("user not logged in")
	}

	user := postgres.GetUserByToken(authorization, &ctx)

	if user == nil {
		return nil, errors.New("user not logged in")
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
		user := postgres.GetUserByToken(authorization, &ctx)

		if user != nil {
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

	if user.Has(auth.RoleApproveMods, &ctx) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canApproveVersions(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(auth.RoleApproveVersions, &ctx) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditUsers(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(auth.RoleEditUsers, &ctx) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditSMLVersions(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(auth.RoleEditSMLVersions, &ctx) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditBootstrapVersions(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(auth.RoleEditBootstrapVersions, &ctx) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditAnnouncements(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(auth.RoleEditAnnouncements, &ctx) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}
