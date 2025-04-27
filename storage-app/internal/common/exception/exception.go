package exception

import (
	"errors"
	"net/http"
)

var (
	ErrTokenExpire                 = errors.New("token expire")
	ErrInvalidToken                = errors.New("invalid token")
	ErrStorageSizeExceedLimitSize  = errors.New("storage size exceed limit size")
	ErrStorageSizeLessThanZero     = errors.New("storage size cannot less than 0")
	ErrMissStorageDetail           = errors.New("storage detail is required")
	ErrInvalidObjectId             = errors.New("object id is invalid")
	ErrParenFileNotExist           = errors.New("parent file not exist")
	ErrTagNotExist                 = errors.New("tag not exist")
	ErrFileNotFound                = errors.New("file not found")
	ErrFileIsFolder                = errors.New("file is folder")
	ErrUnAllowedSortField          = errors.New("unallow sort field")
	ErrUserNotFound                = errors.New("user not found")
	ErrUnSupportOutputImageType    = errors.New("unsupport output image type")
	ErrUserForbiddenFile           = errors.New("user don't have permission access this file")
	ErrEmailExist                  = errors.New("email exist")
	ErrInvalidConfirmCode          = errors.New("invalid confirm code")
	ErrUserPedingSignUpConfirm     = errors.New("user is pending for sign up confirm")
	ErrMismatchPassword            = errors.New("password mismatch")
	ErrWrongPasswordOrEmail        = errors.New("wrong password or email")
	ErrFilePermissionNotFound      = errors.New("file permission not found")
	ErrUnauthorized                = errors.New("unauthorized access")
	ErrPermissionDenied            = errors.New("permission denied")
	ErrMissPermission              = errors.New("miss permission type")
	ErrFileIsNotInBin              = errors.New("file is not in bin")
	ErrFileIsNotFolder             = errors.New("file is not folder")
	ErrSessionNotFound             = errors.New("session not found")
	ErrFileIsUploading             = errors.New("file is in upload progress")
	ErrFileIsNotUploading          = errors.New("file is not uploading")
	ErrUploadFileLockValueMismatch = errors.New("only client upload file can unlock")
	ErrMimeTypeMismatch            = errors.New("update file must be same mimetype")
)
var errMap = map[error]int{
	ErrTokenExpire:                 http.StatusUnauthorized,
	ErrInvalidToken:                http.StatusUnauthorized,
	ErrStorageSizeExceedLimitSize:  http.StatusBadRequest,
	ErrInvalidObjectId:             http.StatusBadRequest,
	ErrParenFileNotExist:           http.StatusNotFound,
	ErrTagNotExist:                 http.StatusNotFound,
	ErrFileNotFound:                http.StatusNotFound,
	ErrUserNotFound:                http.StatusNotFound,
	ErrFileIsFolder:                http.StatusConflict,
	ErrUnAllowedSortField:          http.StatusForbidden,
	ErrUnSupportOutputImageType:    http.StatusBadRequest,
	ErrUserForbiddenFile:           http.StatusForbidden,
	ErrEmailExist:                  http.StatusNotFound,
	ErrInvalidConfirmCode:          http.StatusUnauthorized,
	ErrUserPedingSignUpConfirm:     http.StatusConflict,
	ErrPermissionDenied:            http.StatusForbidden,
	ErrMismatchPassword:            http.StatusConflict,
	ErrWrongPasswordOrEmail:        http.StatusUnauthorized,
	ErrMissStorageDetail:           http.StatusBadRequest,
	ErrFileIsNotInBin:              http.StatusConflict,
	ErrFileIsNotFolder:             http.StatusConflict,
	ErrSessionNotFound:             http.StatusUnauthorized,
	ErrFileIsUploading:             http.StatusConflict,
	ErrFileIsNotUploading:          http.StatusConflict,
	ErrUploadFileLockValueMismatch: http.StatusForbidden,
	ErrMimeTypeMismatch:            http.StatusConflict,
}

func ErrorStatusMapper(err error) int {
	for e, status := range errMap {
		if errors.Is(err, e) {
			return status
		}
	}

	return http.StatusInternalServerError
}
