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
	ErrUserNotFound               = errors.New("user not found")
	ErrUnSupportOutputImageType   = errors.New("unsupport output image type")
	ErrUserForbiddenFile          = errors.New("user don't have permission access this file")
	ErrEmailExist                 = errors.New("email exist")
	ErrInvalidConfirmCode         = errors.New("invalid confirm code")
	ErrUserPedingSignUpConfirm    = errors.New("user is pending for sign up confirm")
	ErrMismatchPassword           = errors.New("password mismatch")
	ErrWrongPasswordOrEmail       = errors.New("wrong password or email")

	ErrUnauthorized     = errors.New("unauthorized access")
	ErrPermissionDenied = errors.New("permission denied")
)
var errMap = map[error]int{
	ErrTokenExpire:                http.StatusUnauthorized,
	ErrInvalidToken:               http.StatusUnauthorized,
	ErrStorageSizeExceedLimitSize: http.StatusBadRequest,
	ErrInvalidObjectId:            http.StatusBadRequest,
	ErrParenFileNotExist:          http.StatusNotFound,
	ErrTagNotExist:                http.StatusNotFound,
	ErrFileNotFound:               http.StatusNotFound,
	ErrUserNotFound:               http.StatusNotFound,
	ErrFileIsFolder:               http.StatusConflict,
	ErrUnAllowedSortField:         http.StatusForbidden,
	ErrUnSupportOutputImageType:   http.StatusBadRequest,
	ErrUserForbiddenFile:          http.StatusForbidden,
	ErrEmailExist:                 http.StatusNotFound,
	ErrInvalidConfirmCode:         http.StatusUnauthorized,
	ErrUserPedingSignUpConfirm:    http.StatusConflict,
	ErrPermissionDenied:           http.StatusForbidden,
	ErrMismatchPassword:           http.StatusConflict,
}

func ErrorStatusMapper(err error) int {
	for e, status := range errMap {
		if errors.Is(err, e) {
			return status
		}
	}

	return http.StatusInternalServerError
}
