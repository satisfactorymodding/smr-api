package auth

type Role struct {
	ID          string
	Description string
}

type Group struct {
	ID    string
	Name  string
	Roles []*Role
}

var (
	RoleApproveMods = &Role{
		ID:          "1",
		Description: "Allows user to approve mods",
	}
	RoleApproveVersions = &Role{
		ID:          "2",
		Description: "Allows user to approve versions",
	}
	RoleDeleteAnyContent = &Role{
		ID:          "3",
		Description: "Allows user to delete content: mod/version/guide",
	}
	RoleEditAnyContent = &Role{
		ID:          "4",
		Description: "Allows user to edit content: mod/version/guide",
	}
	RoleEditUsers = &Role{
		ID:          "5",
		Description: "Allows user to edit other users",
	}
	RoleEditSMLVersions = &Role{
		ID:          "6",
		Description: "Allows user to sml versions",
	}
	RoleEditBootstrapVersions = &Role{
		ID:          "7",
		Description: "Allows user to bootstrap versions",
	}
	RoleEditAnnouncements = &Role{
		ID:          "8",
		Description: "Allows user to manage announcements",
	}
	RoleEditModTags = &Role{
		ID:          "9",
		Description: "Allows user to manage mod tags",
	}
)

var (
	GroupAdmin = &Group{
		ID:   "1",
		Name: "Admin",
		Roles: []*Role{
			RoleApproveMods,
			RoleApproveVersions,
			RoleDeleteAnyContent,
			RoleEditAnyContent,
			RoleEditUsers,
			RoleEditSMLVersions,
			RoleEditBootstrapVersions,
			RoleEditAnnouncements,
			RoleEditModTags,
		},
	}
	GroupModerator = &Group{
		ID:   "2",
		Name: "Moderator",
		Roles: []*Role{
			RoleApproveMods,
			RoleApproveVersions,
			RoleEditAnnouncements,
			RoleEditModTags,
		},
	}
	GroupSMLDev = &Group{
		ID:   "3",
		Name: "SML Dev",
		Roles: []*Role{
			RoleEditSMLVersions,
		},
	}
	GroupBootstrapDev = &Group{
		ID:   "4",
		Name: "Bootstrap Dev",
		Roles: []*Role{
			RoleEditBootstrapVersions,
		},
	}
)

var idToGroupMapping = make(map[string]*Group)
var roleToGroupMapping = make(map[*Role][]*Group)

func initializePermissions() {
	groups := []*Group{GroupAdmin, GroupModerator, GroupSMLDev, GroupBootstrapDev}
	for _, group := range groups {
		idToGroupMapping[group.ID] = group

		for _, role := range group.Roles {
			roleToGroupMapping[role] = append(roleToGroupMapping[role], group)
		}
	}
}

func GetRoleGroups(role *Role) []*Group {
	if groups, ok := roleToGroupMapping[role]; ok {
		return groups
	}

	return []*Group{}
}

func GetGroupByID(id string) *Group {
	return idToGroupMapping[id]
}
