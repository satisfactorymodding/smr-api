package oauth

import (
	"context"

	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/util"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var ctx = context.Background()

var githubAuth *oauth2.Config
var googleAuth *oauth2.Config
var facebookAuth *oauth2.Config

type Site string

const (
	SiteGithub   Site = "github"
	SiteGoogle   Site = "google"
	SiteFacebook Site = "facebook"
)

type UserData struct {
	Email    string
	Username string
	Avatar   string
	Site     Site
	ID       string
}

func InitializeOAuth() {
	githubAuth = &oauth2.Config{
		ClientID:     viper.GetString("oauth.github.client_id"),
		ClientSecret: viper.GetString("oauth.github.client_secret"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	googleAuth = &oauth2.Config{
		ClientID:     viper.GetString("oauth.google.client_id"),
		ClientSecret: viper.GetString("oauth.google.client_secret"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	facebookAuth = &oauth2.Config{
		ClientID:     viper.GetString("oauth.facebook.client_id"),
		ClientSecret: viper.GetString("oauth.facebook.client_secret"),
		Scopes:       []string{"email"},
		Endpoint:     facebook.Endpoint,
	}
}

func GetOAuthOptions(url string) map[string]string {
	urlParam := oauth2.SetAuthURLParam("redirect_uri", url)
	nonce := util.RandomString(16)

	redis.StoreNonce(nonce, url)

	return map[string]string{
		"github":   githubAuth.AuthCodeURL(nonce, oauth2.AccessTypeOffline, urlParam),
		"google":   googleAuth.AuthCodeURL(nonce, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("redirect_uri", viper.GetString("frontend.url"))),
		"facebook": facebookAuth.AuthCodeURL(nonce, urlParam),
	}
}
