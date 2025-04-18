package router

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/libs/pkg/handler"
)

// HandleFunc defines the function signature for handling messages.
type HandleFunc func(context.Context, *sarama.ConsumerMessage) error

// MiddlewareFunc defines the function signature for middleware.
type MiddlewareFunc func(HandleFunc) HandleFunc

type MessageRouter interface {
	Run(context.Context, chan struct{})
	Route(context.Context, *sarama.ConsumerMessage)
	Register(string, HandleFunc, ...MiddlewareFunc)
	RegisterGlobal(...MiddlewareFunc)
}

type MessageRouterImpl struct {
	mapHandlers      map[string]HandleFunc
	mapChs           map[string]chan *sarama.ConsumerMessage
	globalMiddleware []MiddlewareFunc
	errHandler       handler.ErrorHandler
	wg               *sync.WaitGroup
}

// RegisterGlobal implements MessageRouter.
func (m *MessageRouterImpl) RegisterGlobal(middlewares ...MiddlewareFunc) {
	m.globalMiddleware = append(middlewares, m.globalMiddleware...)
}

// NewMessageRouter initializes a new MessageRouterImpl.
func NewMessageRouter(errHandler handler.ErrorHandler) MessageRouter {
	return &MessageRouterImpl{
		mapHandlers:      make(map[string]HandleFunc),
		globalMiddleware: make([]MiddlewareFunc, 0),
		mapChs:           make(map[string]chan *sarama.ConsumerMessage),
		errHandler:       errHandler,
		wg:               &sync.WaitGroup{},
	}
}

// Route routes the message to the appropriate handler, applying middleware.
func (m *MessageRouterImpl) Route(ctx context.Context, msg *sarama.ConsumerMessage) {
	if mapCh, exists := m.mapChs[msg.Topic]; exists {
		mapCh <- msg
	}
	m.errHandler.HandleError(ctx, errors.New("handle function not found for topic: "+msg.Topic))
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
	m.mapChs[topic] = make(chan *sarama.ConsumerMessage, 100)
}

func (m *MessageRouterImpl) Run(ctx context.Context, doneCh chan struct{}) {
	for topic, ch := range m.mapChs {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			for {
				select {
				case <-ctx.Done():
					break
				case msg := <-ch:
					if err := m.mapHandlers[topic](context.Background(), msg); err != nil {
						m.errHandler.HandleError(ctx, err)
					}
					fmt.Println("Handle message success")
				}
			}
		}()
	}
	doneCh <- struct{}{}
	m.wg.Wait()
}
