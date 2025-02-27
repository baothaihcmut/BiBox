package services

import (
	"bytes"
	"context"
	"html/template"

	"github.com/baothaihcmut/BiBox/storage-app-email/internal/config"
	"gopkg.in/gomail.v2"
)

type GmailService interface {
	SendMail(_ context.Context, arg SendMailArg) error
}

type GmailServiceImpl struct {
	dialer     *gomail.Dialer
	mailConfig *config.EmailConfig
}

type SendMailArg struct {
	Subject  string
	To       string
	Template string
	Data     any
}

func (g *GmailServiceImpl) SendMail(_ context.Context, arg SendMailArg) error {
	tmpl, err := template.ParseFiles("templates/confirm_signup.html")
	if err != nil {
		return err
	}
	var body bytes.Buffer
	err = tmpl.Execute(&body, arg.Data)
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", g.mailConfig.Username)
	m.SetHeader("To", arg.To)
	m.SetHeader("Subject", arg.Subject)
	m.SetBody("text/html", body.String())
	if err := g.dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
func NewGmailService(dialer *gomail.Dialer, mailConfig *config.EmailConfig) GmailService {
	return &GmailServiceImpl{
		dialer:     dialer,
		mailConfig: mailConfig,
	}
}
