package router

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
)

type HandleFunc func(context.Context, *sarama.ConsumerMessage) error
type MessageRouter interface {
	Route(context.Context, *sarama.ConsumerMessage) error
	Register(string, HandleFunc)
}

type MessageRouterImpl struct {
	mapHandlers map[string]HandleFunc
}

func (m *MessageRouterImpl) Route(ctx context.Context, msg *sarama.ConsumerMessage) error {
	if handleFunc, exist := m.mapHandlers[msg.Topic]; exist {
		err := handleFunc(ctx, msg)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Hanle function not found")
	}
	return nil
}
