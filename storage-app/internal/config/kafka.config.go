package config

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	// User     string   `mapstructure:"user"`
	// Password string   `mapstructure:"password"`
	MaxRetry int `mapstructure:"max_retry"`
}
