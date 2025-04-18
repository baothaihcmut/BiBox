package handler

import (
	"context"

	"github.com/baothaihcmut/BiBox/libs/pkg/constant"
	"github.com/sirupsen/logrus"
)

type ErrorHandler interface {
	HandleError(context.Context, error)
	Run(context.Context)
}

type errInfo struct {
	RequestId string
	Error     error
}

type ErrorHandlerImpl struct {
	errCh  chan *errInfo
	logger *logrus.Logger
}

func NewErrorHandler(logger *logrus.Logger) *ErrorHandlerImpl {
	return &ErrorHandlerImpl{
		errCh:  make(chan *errInfo, 100),
		logger: logger,
	}
}
func (e *ErrorHandlerImpl) HandleError(ctx context.Context, err error) {
	requestId, ok := ctx.Value(constant.RequesIdContext).(string)
	errInfo := &errInfo{
		Error: err,
	}
	if ok {
		errInfo.RequestId = requestId
	}
	e.errCh <- errInfo
}

func (e *ErrorHandlerImpl) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		case err := <-e.errCh:
			e.logger.Error("Error:", err)
		}
	}
}
