package presenters

type GrantFilePermissionInput struct {
	FileID           string `json:"file_id"`
	UserID           string `json:"user_id"`
	PermissionType   string `json:"permission_type"`
	AccessSecureFile bool   `json:"access_secure_file"`
}

type GrantFilePermissionOutput struct {
	FileID           string `json:"file_id"`
	UserID           string `json:"user_id"`
	PermissionType   string `json:"permission_type"`
	AccessSecureFile bool   `json:"access_secure_file"`
}
