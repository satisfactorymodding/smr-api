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
	"github.com/satisfactorymodding/smr-api/db/schema"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/conv"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/guide"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/generated/ent/usergroup"
	"github.com/satisfactorymodding/smr-api/generated/ent/usermod"
	"github.com/satisfactorymodding/smr-api/generated/ent/usersession"
	"github.com/satisfactorymodding/smr-api/storage"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/util/converter"
)

func (r *mutationResolver) UpdateUser(ctx context.Context, userID string, input generated.UpdateUser) (*generated.User, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "updateUser")
	defer wrapper.end()

	u, err := db.From(ctx).User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	update := u.Update()
	if input.Avatar != nil {
		file, err := io.ReadAll(input.Avatar.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read avatar file: %w", err)
		}

		avatarData, err := converter.ConvertAnyImageToWebp(ctx, file)
		if err != nil {
			return nil, err
		}

		success, avatarKey := storage.UploadUserAvatar(ctx, u.ID, bytes.NewReader(avatarData))
		if success {
			update = update.SetAvatar(storage.GenerateDownloadLink(avatarKey))
		}
	}

	if input.Groups != nil {
		currentGroups, err := db.From(ctx).UserGroup.
			Query().
			Where(usergroup.UserID(u.ID)).
			All(schema.SkipSoftDelete(ctx))
		if err != nil {
			return nil, err
		}

		err = db.Tx(ctx, func(ctx context.Context, tx *ent.Tx) error {
			for _, group := range input.Groups {
				if auth.GetGroupByID(group) == nil {
					continue
				}

				found := false
				var deleted *ent.UserGroup
				for _, currentGroup := range currentGroups {
					if group == currentGroup.GroupID {
						found = true
						if currentGroup.DeletedAt.IsZero() {
							deleted = currentGroup
						}
						break
					}
				}

				if !found {
					if err := tx.UserGroup.
						Create().
						SetUserID(u.ID).
						SetGroupID(group).
						Exec(ctx); err != nil {
						return err
					}
				} else if deleted != nil {
					if err := deleted.Update().ClearDeletedAt().Exec(schema.SkipSoftDelete(ctx)); err != nil {
						return err
					}
				}
			}

			for _, currentGroup := range currentGroups {
				found := false
				for _, group := range input.Groups {
					if group == currentGroup.GroupID {
						found = true
						break
					}
				}

				if !found {
					if _, err := tx.UserGroup.Delete().Where(
						usergroup.UserID(u.ID),
						usergroup.GroupID(currentGroup.GroupID),
					).Exec(ctx); err != nil {
						return err
					}
				}
			}

			return nil
		}, nil)
		if err != nil {
			return nil, err
		}
	}

	if input.Username != nil {
		if len(*input.Username) < 3 {
			return nil, errors.New("username must be at least 3 characters long")
		}

		newUsername := *input.Username
		if len(newUsername) > 32 {
			newUsername = newUsername[:32]
		}

		update = update.SetUsername(newUsername)
	}

	result, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.UserImpl)(nil).Convert(result), nil
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	wrapper, ctx := WrapMutationTrace(ctx, "logout")
	defer wrapper.end()

	header := ctx.Value(util.ContextHeader{}).(http.Header)

	// TODO Archive old deleted sessions to cold storage
	if _, err := db.From(ctx).UserSession.Delete().Where(usersession.Token(header.Get("Authorization"))).Exec(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) GetMe(ctx context.Context) (*generated.User, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "getMe")
	defer wrapper.end()

	result, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.UserImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetUser(ctx context.Context, userID string) (*generated.User, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "getUser")
	defer wrapper.end()

	result, err := db.From(ctx).User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	return (*conv.UserImpl)(nil).Convert(result), nil
}

func (r *queryResolver) GetUsers(ctx context.Context, userIds []string) ([]*generated.User, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "getUsers")
	defer wrapper.end()

	result, err := db.From(ctx).User.Query().Where(user.IDIn(userIds...)).All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.UserImpl)(nil).ConvertSlice(result), nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) Mods(ctx context.Context, obj *generated.User) ([]*generated.UserMod, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "User.mods")
	defer wrapper.end()

	result, err := db.From(ctx).UserMod.
		Query().
		Where(usermod.UserID(obj.ID)).
		Where(usermod.HasModWith(mod.DeletedAtIsNil())).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.UserModImpl)(nil).ConvertSlice(result), nil
}

func (r *userResolver) Guides(ctx context.Context, obj *generated.User) ([]*generated.Guide, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "User.guides")
	defer wrapper.end()

	result, err := db.From(ctx).Guide.Query().Where(guide.UserID(obj.ID)).WithTags().All(ctx)
	if err != nil {
		return nil, err
	}

	return (*conv.GuideImpl)(nil).ConvertSlice(result), nil
}

func (r *userResolver) Groups(ctx context.Context, _ *generated.User) ([]*generated.Group, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "User.guides")
	defer wrapper.end()

	u, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := u.QueryGroups().All(ctx)
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

	u, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := u.QueryGroups().All(ctx)
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

	result, err := dataloader.For(ctx).UserByID.Load(ctx, obj.UserID)()
	if err != nil {
		return nil, err
	}

	return (*conv.UserImpl)(nil).Convert(result), nil
}

func (r *userModResolver) Mod(ctx context.Context, obj *generated.UserMod) (*generated.Mod, error) {
	wrapper, ctx := WrapQueryTrace(ctx, "UserMod.mod")
	defer wrapper.end()

	result, err := db.From(ctx).Mod.Get(ctx, obj.ModID)
	if err != nil {
		return nil, err
	}

	return (*conv.ModImpl)(nil).Convert(result), nil
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

	u, err := db.UserFromGQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, errors.New("user not logged in")
	}

	rawResult := string(nonceString) + "&username=" + u.Username + "&email=" + url.QueryEscape(u.Email) + "&external_id=" + u.ID
	encodedResult := base64.StdEncoding.EncodeToString([]byte(rawResult))
	escapedResult := url.QueryEscape(encodedResult)

	h.Reset()
	h.Write([]byte(encodedResult))
	resultSig := hex.EncodeToString(h.Sum(nil))

	result := viper.GetString("discourse.url") + "?sso=" + escapedResult + "&sig=" + resultSig

	return &result, nil
}
