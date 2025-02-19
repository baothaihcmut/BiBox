package models

type FilePermission struct {
	FileID         string `bson:"file_id,omitempty"`
	UserID         string `bson:"user_id,omitempty"`
	PermissionType string `bson:"permission_type,omitempty"` // e.g., "read", "write", "admin"
	AccessSecure   bool   `bson:"access_secure_file"`
}
