package initialize

import (
	"github.com/baothaihcmut/Storage-app/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func InitializeOauth2(cfg *config.Oauth2Config) *oauth2.Config {
	return &oauth2.Config{
		ClientSecret: cfg.Secret,
		ClientID:     cfg.ClientId,
		RedirectURL:  cfg.RedirectURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
