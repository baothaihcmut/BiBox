package config

type RedisConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Username string `mapstructure:"user_name"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}
