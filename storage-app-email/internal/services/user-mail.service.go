package services

import (
	"github.com/baothaihcmut/BiBox/libs/pkg/events/users"
)

type UserMailService interface {
	SendMailConfirmSignUp(users.UserSignUpEvent)
}
