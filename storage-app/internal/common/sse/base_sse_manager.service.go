package sse

import (
	"context"
	"sync"

	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/cache"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
)

type SSEClient[T any] struct {
	msgCh     chan *T
	sessionId string
}
type SSEManager[T any] struct {
	lock    sync.RWMutex
	clients []*SSEClient[T]
}

type SSEManagerService[T any] struct {
	clientConns  map[string]*SSEManager[T]
	cacheService cache.CacheService
	logger       logger.Logger
}

func NewNotificationSSEManagerService[T any](
	cacheService cache.CacheService,
	logger logger.Logger,
) *SSEManagerService[T] {
	return &SSEManagerService[T]{
		clientConns:  make(map[string]*SSEManager[T]),
		cacheService: cacheService,
		logger:       logger,
	}
}

func (n *SSEManagerService[T]) Connect(ctx context.Context, cacheKey string, sessionId string) (<-chan *T, string, error) {
	key, err := n.cacheService.GetString(ctx, cacheKey+sessionId)
	if err != nil {
		return nil, "", err
	}
	if key == nil {
		return nil, "", exception.ErrSessionNotFound
	}
	if _, exist := n.clientConns[*key]; !exist {
		n.clientConns[*key] = &SSEManager[T]{
			lock:    sync.RWMutex{},
			clients: make([]*SSEClient[T], 0),
		}
	}
	newClient := &SSEClient[T]{
		msgCh:     make(chan *T, 100),
		sessionId: sessionId,
	}
	n.clientConns[*key].clients = append(n.clientConns[*key].clients, newClient)
	return newClient.msgCh, *key, nil
}

func (n *SSEManagerService[T]) SendMessage(ctx context.Context, id string, msg *T) error {
	clientsWithLock, exist := n.clientConns[id]
	if !exist {
		return nil
	}

	clientsWithLock.lock.RLock()
	defer clientsWithLock.lock.RUnlock()

	for _, client := range clientsWithLock.clients {
		select {
		case client.msgCh <- msg:
		default:
			n.logger.Warnf(ctx, nil, "Client channel full for session: %s", client.sessionId)
		}
	}

	return nil
}

func (n *SSEManagerService[T]) ClearClosedSession(ctx context.Context, cacheKey string) {
	wg := &sync.WaitGroup{}

	for _, clients := range n.clientConns {
		clients := clients

		activeClients := make([]*SSEClient[T], 0, len(clients.clients))

		for _, client := range clients.clients {
			wg.Add(1)
			lockArray := sync.RWMutex{}
			go func(client *SSEClient[T]) {
				defer wg.Done()

				userId, err := n.cacheService.GetString(ctx, cacheKey+client.sessionId)
				if err != nil {
					n.logger.Errorf(ctx, nil, "Error checking session: %v", err)
					lockArray.Lock()
					activeClients = append(activeClients, client)
					lockArray.Unlock()
					return
				}

				if userId != nil {
					lockArray.Lock()
					activeClients = append(activeClients, client)
					lockArray.Unlock()
				} else {
					close(client.msgCh)
				}
			}(client)
		}

		go func() {
			wg.Wait()
			clients.lock.Lock()
			clients.clients = activeClients
			clients.lock.Unlock()
		}()
	}
}
func (n *SSEManagerService[T]) Disconnect(ctx context.Context, notificationId string, sessionId string) error {
	clientsWithLock, exist := n.clientConns[notificationId]
	if !exist {
		return nil
	}

	clientsWithLock.lock.RLock()
	defer clientsWithLock.lock.RUnlock()
	for idx, client := range clientsWithLock.clients {
		if client.sessionId == sessionId {
			clientsWithLock.clients = append(clientsWithLock.clients[0:idx], clientsWithLock.clients[idx+1:]...)
			close(client.msgCh)
		}
	}
	return nil
}
