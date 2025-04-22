package impl

import (
	"context"

	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/cache"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/sse"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services"
)

const fileUploadCacheKey = "session:file_upload_progress"

type FileUploadProgressSSEManagerServiceImpl struct {
	*sse.SSEManagerService[response.FileUploadProgressOuput]
	uploadProgressService services.FileUploadProgressService
}

func NewFileUploadProgressSSEManagerService(
	uploadProgressService services.FileUploadProgressService,
	cacheService cache.CacheService,
	logger logger.Logger,
) *FileUploadProgressSSEManagerServiceImpl {
	return &FileUploadProgressSSEManagerServiceImpl{
		SSEManagerService:     sse.NewNotificationSSEManagerService[response.FileUploadProgressOuput](cacheService, logger),
		uploadProgressService: uploadProgressService,
	}
}

func (f *FileUploadProgressSSEManagerServiceImpl) Connect(ctx context.Context, sessionId string) (<-chan *response.FileUploadProgressOuput, string, error) {
	return f.SSEManagerService.Connect(ctx, fileUploadCacheKey, sessionId)
}

func (f *FileUploadProgressSSEManagerServiceImpl) ClearClosedSession(ctx context.Context) {
	f.SSEManagerService.ClearClosedSession(ctx, fileUploadCacheKey)
}
func (f *FileUploadProgressSSEManagerServiceImpl) SendUploadProgressMessage(ctx context.Context, fileId string, e *services.FileUploadProgress) error {
	msg := &response.FileUploadProgressOuput{
		UploadSpeed: e.UploadSpeed,
		Percent:     e.Percent,
		TotalSize:   e.TotalSize,
	}
	return f.SSEManagerService.SendMessage(ctx, fileId, msg)
}

func (f *FileUploadProgressSSEManagerServiceImpl) Disconnect(ctx context.Context, fileId string, sessionId string) error {
	return f.SSEManagerService.Disconnect(ctx, fileId, sessionId)
}
