package config

type TokenConfig struct {
	Secret string `mapstructure:"secret"`
	Age    int    `mapstructure:"age"`
}

type JwtConfig struct {
	AccessToken  TokenConfig `mapstructure:"access_token"`
	RefreshToken TokenConfig `mapstructure:"refresh_token"`
}
