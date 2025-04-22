package impl

import (
	"context"
	"time"

	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/cache"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services"
)

const fileUploadProgressCacheKey = "file:upload_progress:"
const fileUploadProgessChanel = "file:upload_progress_channel"

type FileUploadProgressPayload struct {
	FileId      string  `json:"file_id"`
	UploadSpeed int     `json:"upload_speed"`
	Percent     float32 `json:"percent"`
	TotalSize   int     `json:"total_size"`
}

type FileUploadProgressService struct {
	cacherService cache.CacheService
	logger        logger.Logger
}

func NewFileUploadProgressService(
	cacheService cache.CacheService,
	logger logger.Logger,
) *FileUploadProgressService {
	return &FileUploadProgressService{
		cacherService: cacheService,
		logger:        logger,
	}
}

func (f *FileUploadProgressService) StartUpload(ctx context.Context, fileId string, totalSize int) error {
	uploadProgress := FileUploadProgressPayload{
		FileId:      fileId,
		UploadSpeed: 0,
		Percent:     0,
		TotalSize:   totalSize,
	}
	if err := f.cacherService.SetValue(ctx, fileUploadProgressCacheKey+fileId, uploadProgress, 100*time.Hour); err != nil {
		return err
	}
	return nil
}

func (f *FileUploadProgressService) DoneUpload(ctx context.Context, fileId string) error {
	if err := f.cacherService.Remove(ctx, fileUploadProgressCacheKey+fileId); err != nil {
		return err
	}

	return nil
}

func (f *FileUploadProgressService) GetUploadUploadProgress(ctx context.Context, fileId string) (*services.FileUploadProgress, error) {
	var uploadProgress FileUploadProgressPayload
	if err := f.cacherService.GetValue(ctx, fileUploadProgressCacheKey+fileId, &uploadProgress); err != nil {
		return nil, err
	}
	return &services.FileUploadProgress{
		TotalSize:   uploadProgress.TotalSize,
		Percent:     uploadProgress.Percent,
		UploadSpeed: uploadProgress.UploadSpeed,
	}, nil
}

func (f *FileUploadProgressService) UpdateUploadProgress(ctx context.Context, fileId string, progress *services.FileUploadProgress) error {
	payload := FileUploadProgressPayload{
		FileId:      fileId,
		UploadSpeed: progress.UploadSpeed,
		Percent:     progress.Percent,
		TotalSize:   progress.TotalSize,
	}
	if err := f.cacherService.SetValue(ctx, fileUploadProgressCacheKey+fileId, payload, 100*time.Hour); err != nil {
		return err
	}
	if err := f.cacherService.PublishMessage(ctx, fileUploadProgessChanel, payload); err != nil {
		return err
	}
	return nil
}

func (f *FileUploadProgressService) IsFileUploading(ctx context.Context, fileId string) (bool, error) {
	return f.cacherService.ExistByKey(ctx, fileUploadCacheKey+fileId)
}
