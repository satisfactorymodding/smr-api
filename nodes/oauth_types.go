package nodes

import (
	"github.com/satisfactorymodding/smr-api/db/postgres"
)

type UserSession struct {
	Token string `json:"token"`
}

func SessionToSession(session *postgres.UserSession) *UserSession {
	return &UserSession{
		Token: session.Token,
	}
}
