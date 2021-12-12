package postgres

import (
	"context"

	"github.com/satisfactorymodding/smr-api/auth"
)

func (user User) Has(ctx context.Context, role *auth.Role) bool {
	groups := auth.GetRoleGroups(role)
	groupIds := make([]string, len(groups))
	for i, group := range groups {
		groupIds[i] = group.ID
	}

	var group UserGroup
	DBCtx(ctx).Where("user_id = ?", user.ID).Where("group_id IN (?)", groupIds).First(&group)

	return group.UserID == user.ID
}

func (user User) GetRoles(ctx context.Context) map[*auth.Role]bool {
	var groups []UserGroup
	DBCtx(ctx).Where("user_id = ?", user.ID).Find(&groups)

	roles := make(map[*auth.Role]bool)

	for _, group := range groups {
		gr := auth.GetGroupByID(group.GroupID)
		for _, role := range gr.Roles {
			roles[role] = true
		}
	}

	return roles
}

func (user User) GetGroups(ctx context.Context) []*auth.Group {
	var groups []UserGroup
	DBCtx(ctx).Where("user_id = ?", user.ID).Find(&groups)

	mappedGroups := make([]*auth.Group, len(groups))

	for i, group := range groups {
		mappedGroups[i] = auth.GetGroupByID(group.GroupID)
	}

	return mappedGroups
}

func (user User) SetGroups(ctx context.Context, groups []string) {
	var currentGroups []UserGroup
	DBCtx(ctx).Unscoped().Where("user_id = ?", user.ID).Find(&currentGroups)

	tx := DBCtx(ctx).Begin()

	for _, group := range groups {
		if auth.GetGroupByID(group) == nil {
			continue
		}

		found := false
		var deleted *UserGroup
		for _, currentGroup := range currentGroups {
			if group == currentGroup.GroupID {
				found = true
				if currentGroup.DeletedAt.Valid {
					deleted = &currentGroup
				}
				break
			}
		}
		if !found {
			tx.Create(&UserGroup{
				UserID:  user.ID,
				GroupID: group,
			})
		} else if deleted != nil {
			deleted.DeletedAt.Valid = false
			tx.Unscoped().Save(deleted)
		}
	}

	for _, currentGroup := range currentGroups {
		found := false
		for _, group := range groups {
			if group == currentGroup.GroupID {
				found = true
				break
			}
		}
		if !found {
			tx.Delete(&UserGroup{
				UserID:  user.ID,
				GroupID: currentGroup.GroupID,
			})
		}
	}

	tx.Commit()
}
