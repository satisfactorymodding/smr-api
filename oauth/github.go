package oauth

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/satisfactorymodding/smr-api/redis"

	"github.com/pkg/errors"
)

func GithubCallback(code string, state string) (*UserData, error) {
	_, err := redis.GetNonce(state)

	if err != nil {
		return nil, errors.New("login expired")
	}

	token, err := githubAuth.Exchange(ctx, code)

	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange code")
	}

	client := githubAuth.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user data")
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "failed to read user data")
	}

	var userData map[string]interface{}
	err = json.Unmarshal(bytes, &userData)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user data")
	}

	resp, err = client.Get("https://api.github.com/user/emails")

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user emails")
	}

	bytes, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "failed to read user emails")
	}

	var userEmailsData []map[string]interface{}
	err = json.Unmarshal(bytes, &userEmailsData)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user emails")
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
