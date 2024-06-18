package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Vilsol/slox"
	"github.com/machinebox/graphql"
	"github.com/spf13/viper"

	smr "github.com/satisfactorymodding/smr-api/api"
	"github.com/satisfactorymodding/smr-api/auth"
	"github.com/satisfactorymodding/smr-api/db"
	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/util"
	"github.com/satisfactorymodding/smr-api/validation"

	// Import pgx
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setup() (context.Context, *graphql.Client, func()) {
	validation.StaticPath = "../static"

	client := graphql.NewClient("http://localhost:5020/v2/query")

	ctx, _ := smr.Initialize(context.Background())

	redis.FlushRedis()

	var out []struct {
		TableName string
	}

	connection, err := sql.Open("pgx", fmt.Sprintf(
		"sslmode=disable host=%s port=%d user=%s dbname=%s password=%s",
		viper.GetString("database.postgres.host"),
		viper.GetInt("database.postgres.port"),
		viper.GetString("database.postgres.user"),
		viper.GetString("database.postgres.db"),
		viper.GetString("database.postgres.pass"),
	))
	if err != nil {
		panic(err)
	}

	query, err := connection.Query(`SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`)
	if err != nil {
		panic(err)
	}
	defer query.Close()

	for query.Next() {
		row := struct {
			TableName string
		}{}

		err = query.Scan(&row.TableName)
		if err != nil {
			panic(err)
		}

		out = append(out, row)
	}

	for _, name := range out {
		_, err = connection.Exec(`DROP TABLE IF EXISTS ` + name.TableName + ` CASCADE`)
		if err != nil {
			panic(err)
		}
	}

	smr.Migrate(ctx)
	e := smr.Setup(ctx)
	go smr.Serve(e)

	stopChannel := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-stopChannel
		if err := smr.Stop(e); err != nil {
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
