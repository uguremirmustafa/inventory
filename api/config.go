package api

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	googleOauthConfig *oauth2.Config
	oauthStateString  string
}

func NewConfig() *Config {

	googleOauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_ClientID"),
		ClientSecret: os.Getenv("GOOLE_ClientSecret"),
		RedirectURL:  "http://localhost:9000/v1/auth/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	return &Config{
		googleOauthConfig: googleOauthConfig,
		oauthStateString:  "random_google",
	}
}
