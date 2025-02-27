package config

type EmailConfig struct {
	Username string `mapstructure:"user_name"`
	Password string `mapstructure:"password"`
	MailHost string `mapstructure:"mail_host"`
	MailPort int    `mapstructure:"mail_port"`
}
