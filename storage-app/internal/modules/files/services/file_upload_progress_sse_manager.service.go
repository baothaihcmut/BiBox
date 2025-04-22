package services

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
)

type FileUploadProgressSSEManagerService interface {
	Connect(ctx context.Context, sessionId string) (<-chan *response.FileUploadProgressOuput, string, error)
	ClearClosedSession(ctx context.Context)
	SendUploadProgressMessage(ctx context.Context, fileId string, e *FileUploadProgress) error
	Disconnect(ctx context.Context, fileId string, sessionId string) error
}
