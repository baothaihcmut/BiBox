package interactors

import (
	"context"

	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"github.com/baothaihcmut/Storage-app/internal/modules/auth/presenter"
	"github.com/baothaihcmut/Storage-app/internal/modules/auth/services"
	"github.com/baothaihcmut/Storage-app/internal/modules/users/models"
	"github.com/baothaihcmut/Storage-app/internal/modules/users/repositories"
)

type AuthInteractor interface {
	ExchangeToken(context.Context, *presenter.ExchangeTokenInput) (*presenter.ExchangeTokenOutput, error)
}

type AuthInteractorImpl struct {
	oauth2Service  services.Oauth2Service
	jwtService     services.JwtService
	userRepository repositories.UserRepository
	logger         logger.Logger
}

func NewAuthInteractor(oauth2 services.Oauth2Service, userRepo repositories.UserRepository, jwtService services.JwtService, logger logger.Logger) AuthInteractor {
	return &AuthInteractorImpl{
		userRepository: userRepo,
		oauth2Service:  oauth2,
		jwtService:     jwtService,
		logger:         logger,
	}
}
func (a *AuthInteractorImpl) ExchangeToken(ctx context.Context, input *presenter.ExchangeTokenInput) (*presenter.ExchangeTokenOutput, error) {
	//get user info
	userInfo, err := a.oauth2Service.ExchangeToken(ctx, input.AuthCode)
	if err != nil {
		return nil, err
	}
	//check if user exist in system
	user, err := a.userRepository.FindUserByEmail(ctx, userInfo.Email)
	if err != nil {
		return nil, err
	}
	//if user not exist create new user
	if user == nil {
		user := models.NewUser(userInfo.FirstName, userInfo.LastName, userInfo.Email, userInfo.Image)
		err := a.userRepository.CreateUser(ctx, user)
		if err != nil {
			a.logger.Errorf(ctx, map[string]interface{}{
				"email": userInfo.Email,
			}, "Error create new user:", err)
		}
		a.logger.Info(ctx, map[string]interface{}{
			"email":   userInfo.Email,
			"user_id": user.ID.Hex(),
		}, "User created")
	}
	//generate system token
	accessToken, err := a.jwtService.GenerateAccessToken(ctx, user.ID.Hex())
	if err != nil {
		a.logger.Errorf(ctx, map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
		}, "Error generate token: ", err)
		return nil, err
	}
	refreshToken, err := a.jwtService.GenerateRefreshToken(ctx, user.ID.Hex())
	if err != nil {
		a.logger.Errorf(ctx, map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
		}, "Error generate token: ", err)
		return nil, err
	}
	return &presenter.ExchangeTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
