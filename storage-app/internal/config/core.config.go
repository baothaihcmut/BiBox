package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Logger LoggerConfig `mapstructure:"logger"`
	Server ServerConfig `mapstructure:"server"`
	Mongo  MongoConfig  `mapstructure:"mongo"`
	Jwt    JwtConfig    `mapstructure:"jwt"`
	Oauth2 Oauth2Config `mapstructure:"oauth2"`
	S3     S3Config     `mapstructure:"s3"`
	Kafka  KafkaConfig  `mapstructure:"kafka"`
	Redis  RedisConfig  `mapstructure:"redis"`
}

func LoadConfig() (*AppConfig, error) {
	var config AppConfig
	viper.SetConfigFile("config/config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
