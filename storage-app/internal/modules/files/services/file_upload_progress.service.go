package services

import (
	"context"
)

type FileUploadProgress struct {
	UploadSpeed int
	Percent     float32
	TotalSize   int
}

type FileUploadProgressService interface {
	StartUpload(context.Context, string, int) error
	DoneUpload(context.Context, string) error
	GetUploadUploadProgress(context.Context, string) (*FileUploadProgress, error)
	UpdateUploadProgress(context.Context, string, *FileUploadProgress) error
	IsFileUploading(context.Context, string) (bool, error)
}
