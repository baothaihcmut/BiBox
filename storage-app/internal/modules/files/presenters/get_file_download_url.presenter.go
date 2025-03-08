package presenters

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"

type GetFileDownloadUrlInput struct {
	Id      string `uri:"id" binding:"required"`
	Preview bool   `form:"preview"`
}

type GetFileDownloadUrlOutput struct {
	FileName    string         `json:"file_name"`
	Url         string         `json:"url"`
	Expiry      int            `json:"expiry"`
	Method      string         `json:"method"`
	ContentType enums.MimeType `json:"content_type"`
}
