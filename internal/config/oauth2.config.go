package config

type Oauth2Config struct {
	ClientId    string `mapstructure:"client_id"`
	Secret      string `mapstructure:"secret"`
	RedirectURI string `mapstructure:"redirect_uri"`
}
