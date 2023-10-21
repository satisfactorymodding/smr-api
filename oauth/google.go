package oauth

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/satisfactorymodding/smr-api/redis"
)

func GoogleCallback(code string, state string) (*UserData, error) {
	_, err := redis.GetNonce(state)
	if err != nil {
		return nil, errors.New("login expired")
	}

	token, err := googleAuth.Exchange(ctx, code, oauth2.SetAuthURLParam("redirect_uri", viper.GetString("frontend.url")))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	client := googleAuth.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info: %w", err)
	}

	var userData map[string]interface{}
	err = json.Unmarshal(bytes, &userData)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &UserData{
		Email:    userData["email"].(string),
		Username: userData["name"].(string),
		Avatar:   userData["picture"].(string),
		Site:     SiteGoogle,
		ID:       userData["id"].(string),
	}, nil
}
