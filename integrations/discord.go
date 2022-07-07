package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/satisfactorymodding/smr-api/db/postgres"

	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
	"github.com/russross/blackfriday"
	"github.com/spf13/viper"
)

func NewMod(ctx context.Context, mod *postgres.Mod) {
	if mod == nil {
		return
	}

	if mod.Hidden {
		return
	}

	if viper.GetString("discord.webhook_url") == "" {
		return
	}

	user := postgres.GetUserByID(ctx, mod.CreatorID)

	if user == nil {
		return
	}

	payload := map[string]interface{}{
		"username":   mod.Name,
		"avatar_url": mod.Logo,
		"embeds": []interface{}{
			map[string]interface{}{
				"title":       "**" + mod.Name + "**",
				"url":         "https://ficsit.app/mod/" + mod.ID,
				"color":       16750592,
				"description": mod.ShortDescription,
				"fields": []interface{}{
					map[string]interface{}{
						"name":   "Creator",
						"value":  user.Username,
						"inline": true,
					},
				},
			},
		},
	}

	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		log.Err(err).Msg("error marshaling discord webhook")
		return
	}

	req, _ := http.NewRequest("POST", viper.GetString("discord.webhook_url"), bytes.NewReader(payloadJSON))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	_, _ = ioutil.ReadAll(res.Body)
}

func NewVersion(ctx context.Context, version *postgres.Version) {
	log.Info().Str("stack", string(debug.Stack())).Msg("new version discord webhook")

	if version == nil {
		return
	}

	if viper.GetString("discord.webhook_url") == "" {
		return
	}

	mod := postgres.GetModByID(ctx, version.ModID)

	if mod == nil {
		return
	}

	if mod.Hidden {
		return
	}

	description := version.Changelog
	description = strings.Trim(description, "\n ")
	description = string(blackfriday.MarkdownBasic([]byte(description)))
	description = bluemonday.StrictPolicy().Sanitize(description)
	description = html.UnescapeString(description)
	description = strings.Trim(description, "\n ")

	description = strings.Split(description, "\n")[0]
	if len(description) > 400 {
		description = description[:400] + "..."
	}

	payload := map[string]interface{}{
		"username":   mod.Name,
		"avatar_url": mod.Logo,
		"embeds": []interface{}{
			map[string]interface{}{
				"title":       "**" + mod.Name + " v" + version.Version + "**",
				"url":         "https://ficsit.app/mod/" + mod.ID + "/version/" + version.ID,
				"color":       16750592,
				"description": "New Version Available!",
				"fields": []interface{}{
					map[string]interface{}{
						"name":   "Version",
						"value":  version.Version,
						"inline": true,
					},
					map[string]interface{}{
						"name":   "Stability",
						"value":  version.Stability,
						"inline": true,
					},
				},
				"footer": map[string]interface{}{
					"text": description,
				},
				"thumbnail": map[string]interface{}{
					"url": mod.Logo,
				},
			},
		},
	}

	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		log.Err(err).Msg("error marshaling discord webhook")
		return
	}

	req, _ := http.NewRequest("POST", viper.GetString("discord.webhook_url"), bytes.NewReader(payloadJSON))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	_, _ = ioutil.ReadAll(res.Body)
}
