package gql

import (
	"context"
	"net/http"
	"reflect"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"

	"github.com/satisfactorymodding/smr-api/auth"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/util"
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
		CanManageTags:            canManageTags,
		CanEditModCompatibility:  canEditModCompatibility,
	}
}

type Directive struct {
	generated.DirectiveRoot
}

func canEditMod(ctx context.Context, _ interface{}, next graphql.Resolver, field string) (interface{}, error) {
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

func canEditModCompatibility(ctx context.Context, _ interface{}, next graphql.Resolver, field *string) (interface{}, error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleEditAnyModCompatibility) || user.Has(ctx, auth.RoleEditAnyContent) {
		return next(ctx)
	}

	if field == nil {
		return nil, errors.New("user not authorized to perform this action")
	}

	dbMod := postgres.GetModByID(ctx, getArgument(ctx, *field).(string))

	if dbMod == nil {
		return nil, errors.New("mod not found")
	}

	if postgres.UserCanUploadModVersions(ctx, user, dbMod.ID) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditVersion(ctx context.Context, _ interface{}, next graphql.Resolver, field string) (interface{}, error) {
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

func canEditUser(ctx context.Context, obj interface{}, next graphql.Resolver, field string, object bool) (interface{}, error) {
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

func canEditGuide(ctx context.Context, _ interface{}, next graphql.Resolver, field string) (interface{}, error) {
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

func isLoggedIn(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	header := ctx.Value(util.ContextHeader{}).(http.Header)
	authorization := header.Get("Authorization")

	if authorization == "" {
		return nil, errors.New("user not logged in")
	}

	user := postgres.GetUserByToken(ctx, authorization)

	if user == nil {
		return nil, errors.New("user not logged in")
	}

	if user.Banned {
		return nil, errors.New("user banned")
	}

	userCtx := context.WithValue(ctx, postgres.UserKey{}, user)

	return next(userCtx)
}

func isNotLoggedIn(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	header := ctx.Value(util.ContextHeader{}).(http.Header)
	authorization := header.Get("Authorization")

	if authorization != "" {
		user := postgres.GetUserByToken(ctx, authorization)

		if user != nil {
			return nil, errors.New("user is logged in")
		}
	}

	return next(ctx)
}

func getArgument(ctx context.Context, key string) interface{} {
	return graphql.GetFieldContext(ctx).Args[key]
}

func canApproveMods(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleApproveMods) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canApproveVersions(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleApproveVersions) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditUsers(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleEditUsers) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditSMLVersions(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleEditSMLVersions) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditBootstrapVersions(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleEditBootstrapVersions) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditAnnouncements(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleEditAnnouncements) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canManageTags(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user := ctx.Value(postgres.UserKey{}).(*postgres.User)

	if user.Has(ctx, auth.RoleManageTags) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}
