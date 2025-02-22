package initialize

import (
	"github.com/baothaihcmut/Storage-app/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

func InitializeGoogleOauth2(cfg *config.Oauth2Config) *oauth2.Config {
	return &oauth2.Config{
		ClientSecret: cfg.Google.Secret,
		ClientID:     cfg.Google.ClientId,
		RedirectURL:  cfg.Google.RedirectURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func InitializeFacebookOauth2(cfg *config.Oauth2Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.Facebook.ClientId,
		ClientSecret: cfg.Facebook.Secret,
		RedirectURL:  cfg.Facebook.RedirectURI,
		Scopes: []string{
			"email", "public_profile",
		},
		Endpoint: facebook.Endpoint,
	}
}
