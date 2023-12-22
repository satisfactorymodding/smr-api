package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"

	"github.com/satisfactorymodding/smr-api/redis"
)

func GithubCallback(code string, state string) (*UserData, error) {
	_, err := redis.GetNonce(state)
	if err != nil {
		return nil, errors.New("login expired")
	}

	token, err := githubAuth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	client := githubAuth.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user data: %w", err)
	}

	var userData map[string]interface{}
	err = json.Unmarshal(bytes, &userData)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	resp, err = client.Get("https://api.github.com/user/emails")

	if err != nil {
		return nil, fmt.Errorf("failed to get user emails: %w", err)
	}

	bytes, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read user emails: %w", err)
	}

	var userEmailsData []map[string]interface{}
	err = json.Unmarshal(bytes, &userEmailsData)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user emails: %w", err)
	}

	if len(userEmailsData) == 0 {
		return nil, errors.New("no valid email address found")
	}

	var primaryEmail *string

	for _, v := range userEmailsData {
		if _, ok := v["primary"]; ok {
			if v["primary"] == true {
				email := v["email"].(string)
				primaryEmail = &email
			}
		}
	}

	if primaryEmail == nil {
		return nil, errors.New("no valid email address found")
	}

	avatar := ""

	if avatarData, ok := userData["avatar_url"]; ok {
		avatar = avatarData.(string)
	}

	return &UserData{
		Email:    *primaryEmail,
		Username: userData["login"].(string),
		Avatar:   avatar,
		Site:     SiteGithub,
		ID:       strconv.Itoa(int(userData["id"].(float64))),
	}, nil
}
