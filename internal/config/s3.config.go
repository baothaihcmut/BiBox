package config

type S3Config struct {
	Bucket          string `mapstructure:"bucket"`
	StorageProvider string `mapstructure:"storage_provider"`
	AccessKey       string `mapstructure:"access_key"`
	SecretKey       string `mapstructure:"secret_key"`
	Region          string `mapstructure:"region"`
}
