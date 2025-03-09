package router

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
)

// HandleFunc defines the function signature for handling messages.
type HandleFunc func(context.Context, *sarama.ConsumerMessage) error

// MiddlewareFunc defines the function signature for middleware.
type MiddlewareFunc func(HandleFunc) HandleFunc

type MessageRouter interface {
	Route(context.Context, *sarama.ConsumerMessage) error
	Register(string, HandleFunc, ...MiddlewareFunc)
	RegisterGlobal(...MiddlewareFunc)
}

type MessageRouterImpl struct {
	mapHandlers      map[string]HandleFunc
	globalMiddleware []MiddlewareFunc
}

// RegisterGlobal implements MessageRouter.
func (m *MessageRouterImpl) RegisterGlobal(middlewares ...MiddlewareFunc) {
	m.globalMiddleware = append(middlewares, m.globalMiddleware...)
}

// NewMessageRouter initializes a new MessageRouterImpl.
func NewMessageRouter() MessageRouter {
	return &MessageRouterImpl{
		mapHandlers:      make(map[string]HandleFunc),
		globalMiddleware: make([]MiddlewareFunc, 0),
	}
}

// Route routes the message to the appropriate handler, applying middleware.
func (m *MessageRouterImpl) Route(ctx context.Context, msg *sarama.ConsumerMessage) error {
	if handleFunc, exists := m.mapHandlers[msg.Topic]; exists {
		return handleFunc(ctx, msg)
	}
	return errors.New("handle function not found for topic: " + msg.Topic)
}

// Register registers a handler function for a specific topic, applying middleware.
func (m *MessageRouterImpl) Register(topic string, handler HandleFunc, middlewares ...MiddlewareFunc) {
	// Apply middleware in order
	for _, mwGlobal := range m.globalMiddleware {
		handler = mwGlobal(handler)
	}
	for _, mw := range middlewares {
		handler = mw(handler)
	}
	m.mapHandlers[topic] = handler
}
