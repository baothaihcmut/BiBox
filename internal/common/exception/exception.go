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
)

func ErrorStatusMapper(err error) int {
	switch err {
	case ErrTokenExpire, ErrInvalidToken:
		return http.StatusUnauthorized
	case ErrStorageSizeExceedLimitSize, ErrInvalidObjectId:
		return http.StatusBadRequest
	case ErrParenFileNotExist:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
