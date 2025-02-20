package services

import (
	"bytes"
	"context"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/baothaihcmut/Storage-app/internal/common/enums"
	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"github.com/disintegration/imaging"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
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
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		p.logger.Errorf(ctx, nil, "Error read byte from storage object: ", err)
		return nil, err
	}
	pdf, err := model.NewPdfReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		p.logger.Errorf(ctx, nil, "Error extract object to pdf: ", err)
		return nil, err
	}
	page, err := pdf.GetPage(1)
	if err != nil {
		p.logger.Errorf(ctx, nil, "Error get first page of pdf: ", err)
	}
	ex, err := extractor.New(page)
	if err != nil {
		return nil, err
	}
	imges, err := ex.ExtractPageImages(&extractor.ImageExtractOptions{})
	if err != nil {
		return nil, err
	}
	goImg, err := imges.Images[0].Image.ToGoImage()
	if err != nil {
		p.logger.Errorf(ctx, nil, "Error convert to go image: ", err)
	}
	var result bytes.Buffer
	switch outputType {
	case enums.MimePNG:
		err = png.Encode(&result, imaging.Resize(goImg, 800, 0, imaging.Lanczos))
	case enums.MimeJPG:
		err = jpeg.Encode(&result, imaging.Resize(goImg, 800, 0, imaging.Lanczos), &jpeg.Options{Quality: 10})
	}
	if err != nil {
		p.logger.Errorf(ctx, nil, "Error encode pdf to image: ", err)
		return nil, err
	}
	return &result, nil
}

func (f *FirstPageServiceImpl) GetFirstPage(ctx context.Context, object io.ReadCloser, mimeType enums.MimeType, outputType enums.MimeType) (*bytes.Buffer, error) {
	return f.mapResolveFunc[mimeType](ctx, object, outputType)

}
