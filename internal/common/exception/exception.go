package exception

import (
	"errors"
	"net/http"
)

var (
	ErrTokenExpire  = errors.New("token expire")
	ErrInvalidToken = errors.New("invalid token")
)

func ErrorStatusMapper(err error) int {
	switch err {
	case ErrTokenExpire, ErrInvalidToken:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
