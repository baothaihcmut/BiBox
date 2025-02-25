package services

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/cache"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"github.com/google/uuid"
)

type SendMailConfirmArg struct {
	Email     string
	FirstName string
	LastName  string
	Code      string
}

type UserConfirmService interface {
	StoreUserPending(ctx context.Context, user *models.User) (string, error)
}

type UserConfirmServiceImpl struct {
	cacheService cache.CacheService
	logger       logger.Logger
}

func (u *UserConfirmServiceImpl) StoreUserPending(ctx context.Context, user *models.User) (string, error) {
	//generate code
	code := uuid.New().String()
	//store user info to cache
	err := u.cacheService.SetValue(ctx, fmt.Sprintf("user_pending_confirm:%s", code), user, 30*time.Minute)
	if err != nil {
		u.logger.Errorf(ctx, map[string]interface{}{
			"email": user.Email,
		}, "Error store user info pending confirm to cache: ", err)
		return "", err
	}
	//store email for block user register when pendin
	err = u.cacheService.SetString(ctx, fmt.Sprintf("email_pending_confirm:%s", user.Email), "1", 30*time.Minute)
	if err != nil {
		u.logger.Errorf(ctx, map[string]interface{}{
			"email": user.Email,
		}, "Error store user  email pending confirm to cache: ", err)
		return "", err
	}
	return code, nil
}

func (u *UserConfirmServiceImpl) IsUserPedingConfirm(ctx context.Context, email string) (bool, error) {
	userEmail, err := u.cacheService.GetString(ctx, fmt.Sprintf("email_pending_confirm:%s", email))
	if err != nil {
		return false, err
	}
	return userEmail != nil && *userEmail == "1", nil
}
func (u *UserConfirmServiceImpl) GetUserPedingConfirm(ctx context.Context, code string) (*models.User, error) {
	var user models.User
	err := u.cacheService.GetValue(ctx, fmt.Sprintf("user_pending_confirm:%s", code), &user)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(user, models.User{}) {
		return nil, exception.ErrInvalidConfirmCode
	}
	return &user, nil
}

func (u *UserConfirmServiceImpl) SendMailConfirm(ctx context.Context) {

}
