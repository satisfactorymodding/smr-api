package gql

import (
	"context"
	"reflect"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"

	"github.com/satisfactorymodding/smr-api/auth"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/generated"
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
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if db.UserCanUploadModVersions(ctx, user, getArgument(ctx, field).(string)) {
		return next(ctx)
	}

	if db.UserHas(ctx, auth.RoleEditAnyContent, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditModCompatibility(ctx context.Context, _ interface{}, next graphql.Resolver, field *string) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if db.UserHas(ctx, auth.RoleEditAnyModCompatibility, user) || db.UserHas(ctx, auth.RoleEditAnyContent, user) {
		return next(ctx)
	}

	if field == nil {
		return nil, errors.New("user not authorized to perform this action")
	}

	if db.UserCanUploadModVersions(ctx, user, getArgument(ctx, *field).(string)) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditVersion(ctx context.Context, _ interface{}, next graphql.Resolver, field string) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if db.UserCanUploadModVersions(ctx, user, getArgument(ctx, field).(string)) {
		return next(ctx)
	}

	if db.UserHas(ctx, auth.RoleEditAnyContent, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditUser(ctx context.Context, obj interface{}, next graphql.Resolver, field string, object bool) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	var userID string
	if object {
		userID = reflect.Indirect(reflect.ValueOf(obj)).FieldByName(field).String()
	} else {
		userID = getArgument(ctx, field).(string)
	}

	if userID == user.ID {
		return next(ctx)
	}

	if db.UserHas(ctx, auth.RoleEditUsers, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditGuide(ctx context.Context, _ interface{}, next graphql.Resolver, field string) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	g, err := db.From(ctx).Guide.Get(ctx, getArgument(ctx, field).(string))
	if err != nil {
		return nil, err
	}

	if g.UserID == user.ID {
		return next(ctx)
	}

	if db.UserHas(ctx, auth.RoleEditAnyContent, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func isLoggedIn(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if user.Banned {
		return nil, errors.New("user banned")
	}

	return next(ctx)
}

func isNotLoggedIn(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, errors.New("user is logged in")
	}

	return next(ctx)
}

func getArgument(ctx context.Context, key string) interface{} {
	return graphql.GetFieldContext(ctx).Args[key]
}

func canApproveMods(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if db.UserHas(ctx, auth.RoleApproveMods, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canApproveVersions(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if db.UserHas(ctx, auth.RoleApproveVersions, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditUsers(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if db.UserHas(ctx, auth.RoleEditUsers, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditSMLVersions(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if db.UserHas(ctx, auth.RoleEditSMLVersions, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditBootstrapVersions(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if db.UserHas(ctx, auth.RoleEditBootstrapVersions, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canEditAnnouncements(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if db.UserHas(ctx, auth.RoleEditAnnouncements, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}

func canManageTags(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if db.UserHas(ctx, auth.RoleManageTags, user) {
		return next(ctx)
	}

	return nil, errors.New("user not authorized to perform this action")
}
