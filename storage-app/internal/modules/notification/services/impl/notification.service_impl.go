package impl

import (
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/queue"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/repositories"
)

type NotificationServiceImpl struct {
	repo         repositories.NotificationRepo
	queueService queue.QueueService
	logger       logger.Logger
}

func NewNotificationService(
	repo repositories.NotificationRepo,
	queueService queue.QueueService,
	logger logger.Logger,
) *NotificationServiceImpl {
	return &NotificationServiceImpl{
		repo:         repo,
		queueService: queueService,
		logger:       logger,
	}
}
