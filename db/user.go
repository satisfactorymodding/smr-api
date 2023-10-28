package db

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Vilsol/slox"

	"github.com/satisfactorymodding/smr-api/auth"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/usergroup"
	"github.com/satisfactorymodding/smr-api/generated/ent/usermod"
	"github.com/satisfactorymodding/smr-api/generated/ent/usersession"
	"github.com/satisfactorymodding/smr-api/util"
)

func UserHas(ctx context.Context, role *auth.Role, usr *ent.User) bool {
	groups := auth.GetRoleGroups(role)
	groupIds := make([]string, len(groups))
	for i, group := range groups {
		groupIds[i] = group.ID
	}

	exist, err := usr.QueryGroups().Where(usergroup.GroupIDIn(groupIds...)).Exist(ctx)
	if err != nil {
		slox.Error(ctx, "failed retrieving user groups", slog.Any("err", err))
		return false
	}

	return exist
}

func UserFromGQLContext(ctx context.Context) (*ent.User, error) {
	header := ctx.Value(util.ContextHeader{}).(http.Header)
	authorization := header.Get("Authorization")

	if authorization == "" {
		return nil, errors.New("user not logged in")
	}

	user, err := From(ctx).UserSession.Query().Where(usersession.Token(authorization)).QueryUser().First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("user not logged in")
		}

		return nil, err
	}

	return user, nil
}

func UserCanUploadModVersions(ctx context.Context, user *ent.User, modID string) bool {
	if user.Banned {
		return false
	}

	exists, err := user.QueryUserMods().Where(
		usermod.ModID(modID),
		usermod.RoleIn("creator", "editor"),
	).Exist(ctx)
	if err != nil {
		slox.Error(ctx, "failed retrieving user mods", slog.Any("err", err))
		return false
	}

	return exists
}
