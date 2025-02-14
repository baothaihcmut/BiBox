package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Logger LoggerConfig `mapstructure:"logger"`
	Server ServerConfig `mapstructure:"server"`
	Mongo  MongoConfig  `mapstructure:"mongo"`
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
