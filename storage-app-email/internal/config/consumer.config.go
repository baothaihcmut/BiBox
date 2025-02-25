package config

type ConsumerConfig struct {
	ChannelSize      int      `mapstructure:"channel_size"`
	WorkerPoolSize   int      `mapstructre:"worker_pool_size"`
	Brokers          []string `mapstructure:"brokers"`
	Topics           []string `mapstructure:"topics"`
	ConsumberGroupId string   `mapstructure:"consumer_group_id"`
}
