package models

import "github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"

type UserContext struct {
	Id   string
	Role enums.Role
}
