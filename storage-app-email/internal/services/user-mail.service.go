package services

import (
	"context"

	"github.com/baothaihcmut/BiBox/libs/pkg/events/users"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/models"
)

type UserMailService interface {
	SendMailConfirmSignUp(context.Context, *users.UserSignUpEvent) error
}

type UserMailServiceImpl struct {
	gmailService GmailService
}

func (u *UserMailServiceImpl) SendMailConfirmSignUp(ctx context.Context, e *users.UserSignUpEvent) error {
	args := SendMailArg{
		Subject:  "Confirmation Sign up",
		To:       e.Email,
		Template: "templates/users/confirm_signup.html",
		Data: models.ConfirmSignUpModel{
			FirstName:        e.FirstName,
			LastName:         e.LastName,
			ConfirmationLink: e.ConfirmationLink,
		},
	}
	if err := u.gmailService.SendMail(ctx, args); err != nil {
		return err
	}
	return nil
}
func NewUserMailService(gmailService GmailService) UserMailService {
	return &UserMailServiceImpl{
		gmailService: gmailService,
	}
}
