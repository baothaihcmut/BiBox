package config

type MongoConfig struct {
	Uri               string `mapstructure:"uri"`
	MaxPoolSize       int    `mapstructure:"max_pool_size"`
	MinPoolSize       int    `mapstructure:"min_pool_size"`
	ConnectionTimeout int    `mapstructure:"connection_time_out"`
}
