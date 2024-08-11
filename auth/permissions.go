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
	RoleEditSatisfactoryVersions = &Role{
		ID:          "6",
		Description: "Allows user to edit satisfactory versions",
	}
	RoleEditBootstrapVersions = &Role{
		ID:          "7",
		Description: "Allows user to bootstrap versions",
	}
	RoleEditAnnouncements = &Role{
		ID:          "8",
		Description: "Allows user to manage announcements",
	}
	RoleManageTags = &Role{
		ID:          "9",
		Description: "Allows user to manage tags",
	}
	RoleEditAnyModCompatibility = &Role{
		ID:          "10",
		Description: "Allows user to edit any mod's compatibility info",
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
			RoleEditSatisfactoryVersions,
			RoleEditBootstrapVersions,
			RoleEditAnnouncements,
			RoleManageTags,
			RoleEditAnyModCompatibility,
		},
	}
	GroupModerator = &Group{
		ID:   "2",
		Name: "Moderator",
		Roles: []*Role{
			RoleApproveMods,
			RoleApproveVersions,
			RoleEditAnnouncements,
			RoleManageTags,
			RoleEditAnyModCompatibility,
		},
	}
	GroupSMLDev = &Group{
		ID:   "3",
		Name: "SML Dev",
		Roles: []*Role{
			RoleEditSatisfactoryVersions,
		},
	}
	GroupBootstrapDev = &Group{
		ID:   "4",
		Name: "Bootstrap Dev",
		Roles: []*Role{
			RoleEditBootstrapVersions,
		},
	}
	GroupCompatibilityOfficer = &Group{
		ID:   "5",
		Name: "Compatibility Officer",
		Roles: []*Role{
			RoleEditAnyModCompatibility,
		},
	}
)

var (
	idToGroupMapping   = make(map[string]*Group)
	roleToGroupMapping = make(map[*Role][]*Group)
)

func initializePermissions() {
	groups := []*Group{GroupAdmin, GroupModerator, GroupSMLDev, GroupBootstrapDev, GroupCompatibilityOfficer}
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
