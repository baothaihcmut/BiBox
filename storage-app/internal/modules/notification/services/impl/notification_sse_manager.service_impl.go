package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/BiBox/libs/pkg/events/notifications"
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/cache"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/repositories"
)

const notificationSessionKey = "session:notification"

type NoficationClient struct {
	msgCh     chan *response.NotificationOutput
	sessionId string
}
type NotificationsWithLock struct {
	lock    sync.RWMutex
	clients []*NoficationClient
}

type NotificationSSEManagerService struct {
	clientConns      map[string]*NotificationsWithLock
	notificationRepo repositories.NotificationRepo
	cacheService     cache.CacheService
	logger           logger.Logger
}

func NewNotificationSSEManagerService(
	notificationRepo repositories.NotificationRepo,
	cacheService cache.CacheService,
	logger logger.Logger,
) *NotificationSSEManagerService {
	return &NotificationSSEManagerService{
		clientConns:      make(map[string]*NotificationsWithLock),
		notificationRepo: notificationRepo,
		cacheService:     cacheService,
		logger:           logger,
	}
}

func (n *NotificationSSEManagerService) Connect(ctx context.Context, sessionId string) (<-chan *response.NotificationOutput, string, error) {
	userId, err := n.cacheService.GetString(ctx, notificationSessionKey+sessionId)
	if err != nil {
		return nil, "", err
	}
	if userId == nil {
		return nil, "", exception.ErrSessionNotFound
	}
	if _, exist := n.clientConns[*userId]; !exist {
		n.clientConns[*userId] = &NotificationsWithLock{
			lock:    sync.RWMutex{},
			clients: make([]*NoficationClient, 0),
		}
	}
	newClient := &NoficationClient{
		msgCh:     make(chan *response.NotificationOutput, 100),
		sessionId: sessionId,
	}
	n.clientConns[*userId].clients = append(n.clientConns[*userId].clients, newClient)
	return newClient.msgCh, *userId, nil
}

func (n *NotificationSSEManagerService) SendNotificationCreatedEvent(ctx context.Context, e *notifications.NotificationCreatedEvent) error {
	clientsWithLock, exist := n.clientConns[e.UserId]
	if !exist {
		return nil
	}

	clientsWithLock.lock.RLock()
	defer clientsWithLock.lock.RUnlock()

	for _, client := range clientsWithLock.clients {
		select {
		case client.msgCh <- &response.NotificationOutput{
			Id:         e.Id,
			UserId:     e.UserId,
			Type:       enums.NoficationType(e.Type),
			Title:      e.Title,
			Message:    e.Message,
			ActionUrl:  e.ActionUrl,
			Seen:       e.Seen,
			FromUserId: e.FromUserId,
			CreatedAt:  e.CreatedAt,
		}:
		default:
			n.logger.Warnf(ctx, nil, "Client channel full for session: %s", client.sessionId)
		}
	}

	return nil
}

func (n *NotificationSSEManagerService) ClearClosedSession(ctx context.Context) {
	wg := &sync.WaitGroup{}

	for _, clients := range n.clientConns {
		clients := clients

		activeClients := make([]*NoficationClient, 0, len(clients.clients))

		for _, client := range clients.clients {
			wg.Add(1)
			lockArray := sync.RWMutex{}
			go func(client *NoficationClient) {
				defer wg.Done()

				userId, err := n.cacheService.GetString(ctx, notificationSessionKey+client.sessionId)
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
func (n *NotificationSSEManagerService) Disconnect(ctx context.Context, notificationId string, sessionId string) error {
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
