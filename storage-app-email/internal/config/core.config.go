package config

import (
	"github.com/baothaihcmut/BiBox/libs/pkg/consumer"
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/spf13/viper"
)

type CoreConfig struct {
	Consumer consumer.ConsumerConfig `mapstructure:"consumer"`
	Mail     EmailConfig             `mapstructure:"mail"`
	Logger   logger.LoggerConfig     `mapstructure:"logger"`
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
