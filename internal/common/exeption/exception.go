package exception

import "errors"

var (
	ErrTokenExpire  = errors.New("token expire")
	ErrInvalidToken = errors.New("invalid token")
)
