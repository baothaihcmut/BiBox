package handlers

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/file-process/internal/services"
	"github.com/baothaihcmut/BiBox/libs/pkg/constant"
	"github.com/baothaihcmut/BiBox/libs/pkg/events/files"
	"github.com/baothaihcmut/BiBox/libs/pkg/events/users"
	"github.com/baothaihcmut/BiBox/libs/pkg/middlewares"
	"github.com/baothaihcmut/BiBox/libs/pkg/router"
)

type UserHandler interface {
	Init(r router.MessageRouter)
}

type UserHandlerImpl struct {
	service services.FileProcessService
	msgChs  map[string]chan *sarama.ConsumerMessage
}

func (u *UserHandlerImpl) handleFileUploaded(ctx context.Context, msg *sarama.ConsumerMessage) error {
	//extract event
	e, ok := ctx.Value(constant.PayloadContext).(*files.FileUploadedEvent)
	if !ok {
		return errors.New("invalid payload")
	}
	err := u.service.HandleFileUploaded(ctx, e)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserHandlerImpl) Init(r router.MessageRouter) {
	r.Register("file.uploaded", u.handleFileUploaded, middlewares.ExtractEventMiddleware[users.UserSignUpEvent]())
}
func NewUserHandler(service services.FileProcessService) UserHandler {
	return &UserHandlerImpl{
		service: service,
	}
}
