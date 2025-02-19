package services

import (
	"bytes"
	"context"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/baothaihcmut/Storage-app/internal/common/enums"
	"github.com/baothaihcmut/Storage-app/internal/common/exception"
	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"github.com/gen2brain/go-fitz"
)

type FirstPageService interface {
	GetFirstPage(context.Context, io.ReadCloser, enums.MimeType, enums.MimeType) (*bytes.Buffer, error)
}
type FirstPageServiceImpl struct {
	logger         logger.Logger
	mapResolveFunc map[enums.MimeType]getFirstPageFunc
}

type getFirstPageFunc func(ctx context.Context, reader io.ReadCloser, outputType enums.MimeType) (*bytes.Buffer, error)

func NewFileFirstPageService(logger logger.Logger) FirstPageService {
	service := &FirstPageServiceImpl{
		logger: logger,
	}
	service.mapResolveFunc = map[enums.MimeType]getFirstPageFunc{
		enums.MimePDF: service.getFirstPagePDF,
	}
	return service
}

func (p *FirstPageServiceImpl) getFirstPagePDF(ctx context.Context, reader io.ReadCloser, outputType enums.MimeType) (*bytes.Buffer, error) {
	doc, err := fitz.NewFromReader(reader)
	if err != nil {
		p.logger.Errorf(ctx, map[string]interface{}{
			"file_type": "pdf",
		}, "Error get first page of file:", err)
		return nil, err
	}
	img, err := doc.Image(0)
	if err != nil {
		p.logger.Errorf(ctx, map[string]interface{}{
			"file_type": "pdf",
		}, "Error get first page of file:", err)
		return nil, err
	}
	var imgBuffer bytes.Buffer
	switch outputType {
	case enums.MimeJPG:
		err = jpeg.Encode(&imgBuffer, img, &jpeg.Options{Quality: 10})
		break
	case enums.MimePNG:
		err = png.Encode(&imgBuffer, img)
		break
	default:
		return nil, exception.ErrUnSupportOutputImageType
	}
	if err != nil {
		p.logger.Errorf(ctx, map[string]interface{}{
			"ouput_type": outputType,
		}, "Error encode image to output: ", err)
		return nil, err
	}
	return &imgBuffer, nil
}

func (f *FirstPageServiceImpl) GetFirstPage(ctx context.Context, object io.ReadCloser, mimeType enums.MimeType, outputType enums.MimeType) (*bytes.Buffer, error) {
	return f.mapResolveFunc[mimeType](ctx, object, outputType)

}
