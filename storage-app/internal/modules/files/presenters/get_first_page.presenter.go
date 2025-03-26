package presenters

import (
	"bytes"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
)

type GetFirstPageInput struct {
	FileId    string         `uri:"id" validate:"required"`
	OuputType enums.MimeType `form:"out_type" validate:"required"`
}

type GetFirstPageOutput struct {
	Image *bytes.Buffer
}
