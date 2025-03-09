package handlers

import (
	"context"
	"errors"

	"github.com/IBM/sarama"

	"github.com/baothaihcmut/BiBox/libs/pkg/constant"
	"github.com/baothaihcmut/BiBox/libs/pkg/events/users"
	"github.com/baothaihcmut/BiBox/libs/pkg/middlewares"
	"github.com/baothaihcmut/BiBox/libs/pkg/router"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/services"
)

type UserHandler interface {
	Init(r router.MessageRouter)
}

type UserHandlerImpl struct {
	service services.UserMailService
}

func (u *UserHandlerImpl) handleConfirmSignUp(ctx context.Context, msg *sarama.ConsumerMessage) error {
	//extract event
	e, ok := ctx.Value(constant.PayloadContext).(*users.UserSignUpEvent)
	if !ok {
		return errors.New("invalid payload")
	}
	err := u.service.SendMailConfirmSignUp(ctx, e)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserHandlerImpl) Init(r router.MessageRouter) {
	r.Register("user.sign_up", u.handleConfirmSignUp, middlewares.ExtractEventMiddleware[users.UserSignUpEvent]())
}
func NewUserHandler(service services.UserMailService) UserHandler {
	return &UserHandlerImpl{
		service: service,
	}
}
