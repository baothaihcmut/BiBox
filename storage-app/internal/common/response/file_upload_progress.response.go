package response

type FileUploadProgressOuput struct {
	UploadSpeed int     `json:"upload_speed"`
	Percent     float32 `json:"percent"`
	TotalSize   int     `json:"total_size"`
}
