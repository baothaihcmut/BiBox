package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"

type GetFileDownloadUrlInput struct {
	Id string `uri:"id" binding:"required"`
}

type GetFileDownloadUrlOutput struct {
	Url         string         `json:"url"`
	Expiry      int            `json:"expiry"`
	Method      string         `json:"method"`
	ContentType enums.MimeType `json:"content_type"`
}
