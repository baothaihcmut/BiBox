package enums

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
