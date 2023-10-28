package gql

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/smr-api/auth"
	"github.com/satisfactorymodding/smr-api/dataloader"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/util/converter"
)

func (r *mutationResolver) UpdateUser(ctx context.Context, userID string, input generated.UpdateUser) (*generated.User, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "updateUser")
	defer wrapper.end()

	dbUser := postgres.GetUserByID(newCtx, userID)

	if dbUser == nil {
		return nil, errors.New("user not found")
	}

	if input.Avatar != nil {
		file, err := io.ReadAll(input.Avatar.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read avatar file: %w", err)
		}

		avatarData, err := converter.ConvertAnyImageToWebp(ctx, file)
		if err != nil {
			return nil, err
		}

		success, avatarKey := storage.UploadUserAvatar(ctx, dbUser.ID, bytes.NewReader(avatarData))
		if success {
			dbUser.Avatar = storage.GenerateDownloadLink(avatarKey)
		}
	}

	if input.Groups != nil {
		dbUser.SetGroups(newCtx, input.Groups)
	}

	if input.Username != nil {
		if len(*input.Username) < 3 {
			return nil, errors.New("username must be at least 3 characters long")
		}

		dbUser.Username = *input.Username

		if len(dbUser.Username) > 32 {
			dbUser.Username = dbUser.Username[:32]
		}
	}

	postgres.Save(newCtx, &dbUser)

	return DBUserToGenerated(dbUser), nil
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	wrapper, newCtx := WrapMutationTrace(ctx, "logout")
	defer wrapper.end()

	header := ctx.Value(util.ContextHeader{}).(http.Header)
	postgres.LogoutSession(newCtx, header.Get("Authorization"))

	return true, nil
}

func (r *queryResolver) GetMe(ctx context.Context) (*generated.User, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "getMe")
	defer wrapper.end()

	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.UserImpl)(nil).Convert(user), nil
}

func (r *queryResolver) GetUser(ctx context.Context, userID string) (*generated.User, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getUser")
	defer wrapper.end()
	return DBUserToGenerated(postgres.GetUserByID(newCtx, userID)), nil
}

func (r *queryResolver) GetUsers(ctx context.Context, userIds []string) ([]*generated.User, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "getUsers")
	defer wrapper.end()

	users := postgres.GetUsersByID(newCtx, userIds)

	if users == nil {
		return nil, errors.New("users not found")
	}

	converted := make([]*generated.User, len(*users))
	for k, v := range *users {
		converted[k] = DBUserToGenerated(&v)
	}

	return converted, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) Mods(ctx context.Context, obj *generated.User) ([]*generated.UserMod, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "User.mods")
	defer wrapper.end()

	mods := postgres.GetUserMods(newCtx, obj.ID)

	if mods == nil {
		return []*generated.UserMod{}, nil
	}

	converted := make([]*generated.UserMod, len(mods))
	for k, v := range mods {
		converted[k] = &generated.UserMod{
			UserID: v.UserID,
			ModID:  v.ModID,
			Role:   v.Role,
		}
	}

	return converted, nil
}

func (r *userResolver) Guides(ctx context.Context, obj *generated.User) ([]*generated.Guide, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "User.guides")
	defer wrapper.end()

	guides := postgres.GetUserGuides(newCtx, obj.ID)

	if guides == nil {
		return nil, errors.New("guides not found")
	}

	converted := make([]*generated.Guide, len(guides))
	for k, v := range guides {
		converted[k] = DBGuideToGenerated(&v)
	}

	return converted, nil
}

func (r *userResolver) Groups(ctx context.Context, _ *generated.User) ([]*generated.Group, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "User.guides")
	defer wrapper.end()

	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := user.QueryGroups().All(ctx)
	if err != nil {
		return nil, err
	}

	converted := make([]*generated.Group, len(groups))
	for k, v := range groups {
		g := auth.GetGroupByID(v.GroupID)
		converted[k] = &generated.Group{
			ID:   g.ID,
			Name: g.Name,
		}
	}

	return converted, nil
}

func (r *userResolver) Roles(ctx context.Context, _ *generated.User) (*generated.UserRoles, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "User.guides")
	defer wrapper.end()

	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := user.QueryGroups().All(ctx)
	if err != nil {
		return nil, err
	}

	roles := make(map[*auth.Role]bool)

	for _, group := range groups {
		gr := auth.GetGroupByID(group.GroupID)
		for _, role := range gr.Roles {
			roles[role] = true
		}
	}

	userRoles := &generated.UserRoles{}

	if hasRole, ok := roles[auth.RoleApproveMods]; ok && hasRole {
		userRoles.ApproveMods = true
	}

	if hasRole, ok := roles[auth.RoleApproveVersions]; ok && hasRole {
		userRoles.ApproveVersions = true
	}

	if hasRole, ok := roles[auth.RoleDeleteAnyContent]; ok && hasRole {
		userRoles.DeleteContent = true
	}

	if hasRole, ok := roles[auth.RoleEditAnyContent]; ok && hasRole {
		userRoles.EditContent = true
	}

	if hasRole, ok := roles[auth.RoleEditUsers]; ok && hasRole {
		userRoles.EditUsers = true
	}

	if hasRole, ok := roles[auth.RoleEditSMLVersions]; ok && hasRole {
		userRoles.EditSMLVersions = true
	}

	if hasRole, ok := roles[auth.RoleEditBootstrapVersions]; ok && hasRole {
		userRoles.EditBootstrapVersions = true
	}

	if hasRole, ok := roles[auth.RoleEditAnyModCompatibility]; ok && hasRole {
		userRoles.EditAnyModCompatibility = true
	}

	return userRoles, nil
}

type userModResolver struct{ *Resolver }

func (r *userModResolver) User(ctx context.Context, obj *generated.UserMod) (*generated.User, error) {
	wrapper, _ := WrapQueryTrace(ctx, "UserMod.user")
	defer wrapper.end()

	user, err := dataloader.For(ctx).UserByID.Load(obj.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return DBUserToGenerated(user), nil
}

func (r *userModResolver) Mod(ctx context.Context, obj *generated.UserMod) (*generated.Mod, error) {
	wrapper, newCtx := WrapQueryTrace(ctx, "UserMod.mod")
	defer wrapper.end()

	mod := postgres.GetModByID(newCtx, obj.ModID)

	if mod == nil {
		return nil, errors.New("mod not found")
	}

	return DBModToGenerated(mod), nil
}

func (r *mutationResolver) DiscourseSso(ctx context.Context, sso string, sig string) (*string, error) {
	wrapper, _ := WrapMutationTrace(ctx, "discourseSSO")
	defer wrapper.end()

	h := hmac.New(sha256.New, []byte(viper.GetString("discourse.sso_secret")))
	h.Write([]byte(sso))

	if sig != hex.EncodeToString(h.Sum(nil)) {
		return nil, errors.New("invalid signature")
	}

	nonceString, err := base64.StdEncoding.DecodeString(sso)
	if err != nil {
		return nil, fmt.Errorf("failed to decode sso: %w", err)
	}

	user, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not logged in")
	}

	rawResult := string(nonceString) + "&username=" + user.Username + "&email=" + url.QueryEscape(user.Email) + "&external_id=" + user.ID
	encodedResult := base64.StdEncoding.EncodeToString([]byte(rawResult))
	escapedResult := url.QueryEscape(encodedResult)

	h.Reset()
	h.Write([]byte(encodedResult))
	resultSig := hex.EncodeToString(h.Sum(nil))

	result := viper.GetString("discourse.url") + "?sso=" + escapedResult + "&sig=" + resultSig

	return &result, nil
}
