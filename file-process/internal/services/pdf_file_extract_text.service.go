package services

import (
	"bytes"
	"context"

	"github.com/baothaihcmut/BiBox/libs/pkg/events/files"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

type PDFFileExtractTextService struct {
}

func (p *PDFFileExtractTextService) Process(ctx context.Context, file bytes.Buffer) (string, error) {
	reader, err := model.NewPdfReader(bytes.NewReader(file.Bytes()))
	if err != nil {
		return "", err
	}
	numPages, err := reader.GetNumPages()
	if err != nil {
		return "", err
	}

	var text string
	for i := 1; i <= numPages; i++ {
		page, err := reader.GetPage(i)
		if err != nil {
			return "", err
		}
		ex, err := extractor.New(page)
		if err != nil {
			return "", err
		}
		pageText, err := ex.ExtractText()
		if err != nil {
			return "", err
		}
		text += pageText + "\n"
	}
	return text, nil

}
func (d *PDFFileExtractTextService) GetFileType() files.MimeType {
	return files.MimePDF
}
