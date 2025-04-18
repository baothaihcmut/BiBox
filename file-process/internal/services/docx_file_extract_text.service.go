package services

import (
	"bytes"
	"context"

	"github.com/baothaihcmut/BiBox/libs/pkg/events/files"
	"github.com/unidoc/unioffice/document"
)

type DocxExtractFileTextService struct {
}

func (d *DocxExtractFileTextService) GetFileType() files.MimeType {
	return files.MimeDOCX
}

func (d *DocxExtractFileTextService) Process(ctx context.Context, file bytes.Buffer) (string, error) {
	doc, err := document.Read(bytes.NewReader(file.Bytes()), int64(file.Len()))
	if err != nil {
		return "", err
	}
	var text string
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text += run.Text()
		}
		text += "\n"
	}

	return text, nil

}
