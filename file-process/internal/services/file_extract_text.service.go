package services

import (
	"bytes"
	"context"

	"github.com/baothaihcmut/BiBox/libs/pkg/events/files"
)

type FileExtractTextService interface {
	Process(context.Context, bytes.Buffer) (string, error)
	GetFileType() files.MimeType
}
