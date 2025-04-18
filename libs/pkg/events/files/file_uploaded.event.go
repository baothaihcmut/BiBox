package files

import (
	"mime"
	"path/filepath"
)

// MimeType represents supported MIME types as an enum
type MimeType string

const (
	// Images
	MimeJPG  MimeType = "image/jpeg"
	MimePNG  MimeType = "image/png"
	MimeGIF  MimeType = "image/gif"
	MimeBMP  MimeType = "image/bmp"
	MimeWEBP MimeType = "image/webp"
	MimeTIFF MimeType = "image/tiff"
	MimeSVG  MimeType = "image/svg+xml"

	// Documents
	MimePDF  MimeType = "application/pdf"
	MimeDOC  MimeType = "application/msword"
	MimeDOCX MimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	MimeXLS  MimeType = "application/vnd.ms-excel"
	MimeXLSX MimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MimePPT  MimeType = "application/vnd.ms-powerpoint"
	MimePPTX MimeType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	MimeTXT  MimeType = "text/plain"
	MimeCSV  MimeType = "text/csv"
	MimeJSON MimeType = "application/json"
	MimeXML  MimeType = "application/xml"
	MimeYAML MimeType = "application/x-yaml"

	// Audio
	MimeMP3  MimeType = "audio/mpeg"
	MimeWAV  MimeType = "audio/wav"
	MimeOGG  MimeType = "audio/ogg"
	MimeFLAC MimeType = "audio/flac"
	MimeAAC  MimeType = "audio/aac"

	// Video
	MimeMP4  MimeType = "video/mp4"
	MimeWebM MimeType = "video/webm"
	MimeAVI  MimeType = "video/x-msvideo"
	MimeMOV  MimeType = "video/quicktime"
	MimeMKV  MimeType = "video/x-matroska"

	// Archives
	MimeZIP  MimeType = "application/zip"
	MimeRAR  MimeType = "application/vnd.rar"
	Mime7z   MimeType = "application/x-7z-compressed"
	MimeTAR  MimeType = "application/x-tar"
	MimeGZIP MimeType = "application/gzip"

	// Code Files
	MimeHTML   MimeType = "text/html"
	MimeCSS    MimeType = "text/css"
	MimeJS     MimeType = "application/javascript"
	MimeGo     MimeType = "text/x-go"
	MimePython MimeType = "text/x-python"
	MimeShell  MimeType = "application/x-sh"

	// Executables
	MimeEXE MimeType = "application/x-msdownload"
	MimeBIN MimeType = "application/octet-stream"
)

var MimeTypeMap = map[string]MimeType{
	"image/jpeg":    MimeJPG,
	"image/png":     MimePNG,
	"image/gif":     MimeGIF,
	"image/bmp":     MimeBMP,
	"image/webp":    MimeWEBP,
	"image/tiff":    MimeTIFF,
	"image/svg+xml": MimeSVG,

	"application/pdf":    MimePDF,
	"application/msword": MimeDOC,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": MimeDOCX,
	"application/vnd.ms-excel": MimeXLS,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         MimeXLSX,
	"application/vnd.ms-powerpoint":                                             MimePPT,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": MimePPTX,
	"text/plain":         MimeTXT,
	"text/csv":           MimeCSV,
	"application/json":   MimeJSON,
	"application/xml":    MimeXML,
	"application/x-yaml": MimeYAML,

	"audio/mpeg": MimeMP3,
	"audio/wav":  MimeWAV,
	"audio/ogg":  MimeOGG,
	"audio/flac": MimeFLAC,
	"audio/aac":  MimeAAC,

	"video/mp4":        MimeMP4,
	"video/webm":       MimeWebM,
	"video/x-msvideo":  MimeAVI,
	"video/quicktime":  MimeMOV,
	"video/x-matroska": MimeMKV,

	"application/zip":             MimeZIP,
	"application/vnd.rar":         MimeRAR,
	"application/x-7z-compressed": Mime7z,
	"application/x-tar":           MimeTAR,
	"application/gzip":            MimeGZIP,

	"text/html":              MimeHTML,
	"text/css":               MimeCSS,
	"application/javascript": MimeJS,
	"text/x-go":              MimeGo,
	"text/x-python":          MimePython,
	"application/x-sh":       MimeShell,

	"application/x-msdownload": MimeEXE,
	"application/octet-stream": MimeBIN,
}

// Function to check if a MIME type is allowed
func IsAllowedMimeType(mimeType MimeType) bool {
	allowedMimeTypes := map[MimeType]bool{
		MimeJPG:  true,
		MimePNG:  true,
		MimeGIF:  true,
		MimePDF:  true,
		MimeDOCX: true,
		MimeXLSX: true,
		MimeMP3:  true,
		MimeMP4:  true,
		MimeZIP:  true,
	}
	return allowedMimeTypes[mimeType]
}

func MapToMimeType(fileName, fallback string) MimeType {
	ext := filepath.Ext(fileName)
	mimeStr := mime.TypeByExtension(ext)
	if mimeStr == "" {
		mimeStr = fallback
	}
	if mimeType, exist := MimeTypeMap[mimeStr]; exist {
		return mimeType
	} else {
		return MimeBIN
	}
}
func GetMimeTypePointer(t MimeType) *MimeType {
	return &t
}

type FileUploadedEvent struct {
	Id              string   `json:"id"`
	StorageKey      string   `json:"storage_key"`
	StorageProvider string   `json:"storage_provider"`
	Size            int      `json:"size"`
	MimeType        MimeType `json:"mime_type"`
}
