package enums

type FilePermissionType int

const (
	ViewPermission FilePermissionType = iota + 1
	CommentPermission
	EditPermission
)

func GetPermissionTypePointer(p PermissionType) *PermissionType {
	return &p
}
