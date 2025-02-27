package config

type KafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	MaxRetry int      `mapstructure:"max_retry"`
}
