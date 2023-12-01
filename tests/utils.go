package tests

import (
	"context"
	"log/slog"
	"sync"
	"github.com/Vilsol/slox"
	"github.com/machinebox/graphql"

	smr "github.com/satisfactorymodding/smr-api/api"
	"github.com/satisfactorymodding/smr-api/auth"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"
)

func setup() (context.Context, *graphql.Client, func()) {
	validation.StaticPath = "../static"

	client := graphql.NewClient("http://localhost:5020/v2/query")

	ctx := smr.Initialize(context.Background())

	redis.FlushRedis()

	var out []struct {
		TableName string
	}

	// TODO Replace with ENT
	err := postgres.DBCtx(ctx).Raw(`SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`).Scan(&out).Error
	if err != nil {
		panic(err)
	}

	for _, name := range out {
		err := postgres.DBCtx(ctx).Exec(`DROP TABLE IF EXISTS ` + name.TableName + ` CASCADE`).Error
		if err != nil {
			panic(err)
		}
	}

	smr.Migrate(ctx)
	smr.Setup(ctx)
	go smr.Serve()

	stopChannel := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-stopChannel
		if err := smr.Stop(); err != nil {
			panic(err)
		}
	}()

	return ctx, client, func() {
		stopChannel <- true
		wg.Wait()
	}
}

func makeUser(ctx context.Context) (string, string, error) {
	user := db.From(ctx).User.
		Create().
		SetEmail("test_user@ficsit.app").
		SetUsername("test_user").
		SaveX(ctx)

	slox.Info(ctx, "created fake test_user", slog.String("id", user.ID))

	db.From(ctx).UserGroup.Create().SetUser(user).SetGroupID(auth.GroupAdmin.ID).SaveX(ctx)

	slox.Info(ctx, "created user admin group")

	session := db.From(ctx).UserSession.
		Create().
		SetUser(user).
		SetToken(util.GenerateUserToken()).
		SaveX(ctx)

	slox.Info(ctx, "created fake user session", slog.String("token", session.Token))

	return session.Token, user.ID, nil
}

func authRequest(q string, token string) *graphql.Request {
	req := graphql.NewRequest(q)
	req.Header.Set("Authorization", token)
	return req
}
