package main

import (
	"fmt"

	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/config"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/server"
	"gopkg.in/gomail.v2"
)

func main() {
	//cfg
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error load config")
		panic(err)
	}
	//mail
	mailDialer := gomail.NewDialer(cfg.Mail.MailHost, cfg.Mail.MailPort, cfg.Mail.Username, cfg.Mail.Password)
	//logrus
	logrus := logger.InitializeLogger(&cfg.Logger)
	s := server.NewServer(
		mailDialer,
		cfg,
		logrus,
	)
	s.Run()

}
