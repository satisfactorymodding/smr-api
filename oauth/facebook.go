package oauth

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/smr-api/redis"
	"golang.org/x/oauth2"
)

func FacebookCallback(code string, state string) (*UserData, error) {
	redirectURI, err := redis.GetNonce(state)

	if err != nil {
		return nil, errors.New("login expired")
	}

	urlParam := oauth2.SetAuthURLParam("redirect_uri", redirectURI)

	token, err := facebookAuth.Exchange(ctx, code, urlParam)

	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange code")
	}

	client := facebookAuth.Client(ctx, token)

	resp, err := client.Get("https://graph.facebook.com/v5.0/me?fields=email,short_name,id,picture{url}")

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user data")
	}

	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var userData map[string]interface{}
	err = json.Unmarshal(bytes, &userData)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user data")
	}

	return &UserData{
		Email:    userData["email"].(string),
		Username: userData["short_name"].(string),
		Avatar:   userData["picture"].(map[string]interface{})["data"].(map[string]interface{})["url"].(string),
		Site:     SiteFacebook,
		ID:       userData["id"].(string),
	}, nil
}
