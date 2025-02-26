package services

import (
	"github.com/baothaihcmut/BiBox/libs/pkg/events"
)

type UserMailService interface {
	SendMailConfirmSignUp(events.UserSignUpEvent)
}
