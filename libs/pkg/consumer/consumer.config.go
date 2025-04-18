package consumer

type ConsumerConfig struct {
	WorkerPoolSize   int      `mapstructure:"worker_pool_size"`
	Brokers          []string `mapstructure:"brokers"`
	Topics           []string `mapstructure:"topics"`
	ConsumberGroupId string   `mapstructure:"consumer_group_id"`
}
