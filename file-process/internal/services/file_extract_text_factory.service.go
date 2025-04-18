package services

import (
	"context"
	"errors"

	"github.com/baothaihcmut/BiBox/libs/pkg/events/files"
	"github.com/samber/lo"
)

type FileExtractTextFactory interface {
	GetService(context.Context, files.MimeType) (FileExtractTextService, error)
}

type FileExtractTextFactoryImpl struct {
	fileExtractTextServices []FileExtractTextService
}

func NewFileExtractTextFactory(services ...FileExtractTextService) *FileExtractTextFactoryImpl {
	return &FileExtractTextFactoryImpl{
		fileExtractTextServices: services,
	}
}

func (f *FileExtractTextFactoryImpl) GetService(ctx context.Context, mimeType files.MimeType) (FileExtractTextService, error) {
	svc := lo.Filter(f.fileExtractTextServices, func(svc FileExtractTextService, _ int) bool {
		return svc.GetFileType() == mimeType
	})
	if len(svc) == 0 {
		return nil, errors.New("file type not support")
	}
	return svc[0], nil
}
