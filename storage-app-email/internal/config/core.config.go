package config

import "github.com/spf13/viper"

type CoreConfig struct {
	Consumer ConsumerConfig `mapstructure:"consumer"`
	Mail     EmailConfig    `mapstructure:"mail"`
}

func LoadConfig() (*CoreConfig, error) {
	var config CoreConfig
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
