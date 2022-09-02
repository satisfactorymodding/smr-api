package oauth

import (
	"encoding/json"
	"io"

	"github.com/satisfactorymodding/smr-api/redis"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func GoogleCallback(code string, state string) (*UserData, error) {
	_, err := redis.GetNonce(state)

	if err != nil {
		return nil, errors.New("login expired")
	}

	token, err := googleAuth.Exchange(ctx, code, oauth2.SetAuthURLParam("redirect_uri", viper.GetString("frontend.url")))

	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange code")
	}

	client := googleAuth.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user info")
	}

	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "failed to read user info")
	}

	var userData map[string]interface{}
	err = json.Unmarshal(bytes, &userData)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user info")
	}

	return &UserData{
		Email:    userData["email"].(string),
		Username: userData["name"].(string),
		Avatar:   userData["picture"].(string),
		Site:     SiteGoogle,
		ID:       userData["id"].(string),
	}, nil
}
