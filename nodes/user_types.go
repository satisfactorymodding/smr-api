package nodes

import (
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
)

type User struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
}

func UserToPrivateUser(user *postgres.User) *User {
	return &User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
	}
}

type PublicUser struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
}

func UserToPublicUser(user *postgres.User) *PublicUser {
	return &PublicUser{
		ID:        user.ID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
	}
}

type UserMod struct {
	ModID string `json:"mod_id"`
	Role  string `json:"role"`
}

func UserModToUserMod(mod *postgres.UserMod) *UserMod {
	return &UserMod{
		ModID: mod.ModID,
		Role:  mod.Role,
	}
}
