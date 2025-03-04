package enums

type PermissionType int

const (
	PermissionTypeView    PermissionType = iota + 1
	PermissionTypeComment PermissionType = iota + 2
	PermissionEdit        PermissionType = iota + 3
)
