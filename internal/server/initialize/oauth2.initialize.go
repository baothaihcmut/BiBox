package initialize

import (
	"github.com/baothaihcmut/Storage-app/internal/config"
	"golang.org/x/oauth2"
)

func InitializeOauth2(cfg *config.Oauth2ConfigInfo, scopes []string, endpoint oauth2.Endpoint) *oauth2.Config {
	return &oauth2.Config{
		ClientSecret: cfg.Secret,
		ClientID:     cfg.ClientId,
		RedirectURL:  cfg.RedirectURI,
		Scopes:       scopes,
		Endpoint:     endpoint,
	}
}
