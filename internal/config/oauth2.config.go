package config

type Oauth2Config struct {
	Google   Oauth2ConfigInfo `mapstructure:"google"`
	Facebook Oauth2ConfigInfo `mapstructure:"facebook"`
}

type Oauth2ConfigInfo struct {
	ClientId    string `mapstructure:"client_id"`
	Secret      string `mapstructure:"secret"`
	RedirectURI string `mapstructure:"redirect_uri"`
}
