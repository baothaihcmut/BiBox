package services

import (
	"context"

	"github.com/baothaihcmut/BiBox/libs/pkg/events/files"
)

type FileProcessService interface {
	HandleFileUploaded(context.Context, *files.FileUploadedEvent) error
}

type FileProcressServiceImpl struct {
	fileExtractTextFactory FileExtractTextFactory
	indexTextService       IndexTextService
	storageService         StorageService
}

func (f *FileProcressServiceImpl) HandleFileUploaded(ctx context.Context, e *files.FileUploadedEvent) error {
	file, err := f.storageService.GetFile(ctx, e.StorageKey)
	if err != nil {
		return err
	}
	svc, err := f.fileExtractTextFactory.GetService(ctx, e.MimeType)
	if err != nil {
		return err
	}
	content, err := svc.Process(ctx, file)
	if err != nil {
		return err
	}
	if err := f.indexTextService.IndexFile(ctx, e.Id, content); err != nil {
		return err
	}
	return nil

}
