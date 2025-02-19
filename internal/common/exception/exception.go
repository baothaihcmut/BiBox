package exception

import (
	"errors"
	"net/http"
)

var (
	ErrTokenExpire                = errors.New("token expire")
	ErrInvalidToken               = errors.New("invalid token")
	ErrStorageSizeExceedLimitSize = errors.New("storage size exceed limit size")
	ErrStorageSizeLessThanZero    = errors.New("storage size cannot less than 0")
	ErrInvalidObjectId            = errors.New("object id is invalid")
	ErrParenFileNotExist          = errors.New("parent file not exist")
	ErrTagNotExist                = errors.New("tag not exist")
	ErrFileNotFound               = errors.New("file not found")
	ErrFileIsFolder               = errors.New("folder cannot be uploaded")
	ErrUnAllowedSortField         = errors.New("unallow sort field")
)

func ErrorStatusMapper(err error) int {
	switch err {
	case ErrTokenExpire, ErrInvalidToken:
		return http.StatusUnauthorized
	case ErrStorageSizeExceedLimitSize, ErrInvalidObjectId:
		return http.StatusBadRequest
	case ErrParenFileNotExist, ErrTagNotExist, ErrFileNotFound:
		return http.StatusNotFound
	case ErrFileIsFolder:
		return http.StatusConflict
	case ErrUnAllowedSortField:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
