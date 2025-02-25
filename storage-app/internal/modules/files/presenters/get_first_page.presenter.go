package presenters

import (
	"bytes"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
)

type GetFirstPageInput struct {
	FileId    string         `uri:"id" binding:"required oneof=image/jpeg image/png"`
	OuputType enums.MimeType `form:"out_type" binding:"required"`
}

type GetFirstPageOutput struct {
	Image *bytes.Buffer
}
