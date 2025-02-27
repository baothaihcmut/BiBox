package interactors

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/presenter"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/services"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/repositories"
)

type AuthInteractor interface {
	ExchangeToken(context.Context, *presenter.ExchangeTokenInput) (*presenter.ExchangeTokenOutput, error)
	SignUp(context.Context, *presenter.SignUpInput) (*presenter.SignUpOuput, error)
}

type AuthInteractorImpl struct {
	oauth2ServiceFactory services.Oauth2ServiceFactory
	jwtService           services.JwtService
	userRepository       repositories.UserRepository
	userConfirmService   services.UserConfirmService
	logger               logger.Logger
}

func NewAuthInteractor(
	oauth2 services.Oauth2ServiceFactory,
	userRepo repositories.UserRepository,
	jwtService services.JwtService,
	logger logger.Logger,
	userConfirmService services.UserConfirmService,
) AuthInteractor {
	return &AuthInteractorImpl{
		userRepository:       userRepo,
		oauth2ServiceFactory: oauth2,
		jwtService:           jwtService,
		logger:               logger,
		userConfirmService:   userConfirmService,
	}
}
func (a *AuthInteractorImpl) ExchangeToken(ctx context.Context, input *presenter.ExchangeTokenInput) (*presenter.ExchangeTokenOutput, error) {
	//get service
	oauth2Service := a.oauth2ServiceFactory.GetOauth2Service(services.Oauth2ServiceToken(input.Provider))
	//get user info
	userInfo, err := oauth2Service.ExchangeToken(ctx, input.AuthCode)
	if err != nil {
		return nil, err
	}
	//check if user exist in system
	user, err := a.userRepository.FindUserByEmail(ctx, userInfo.GetEmail())
	if err != nil {
		return nil, err
	}
	//if user not exist create new user
	imageUrl := userInfo.GetImage()
	if user == nil {
		user = models.NewUser(
			userInfo.GetFirstName(),
			userInfo.GetLastName(),
			userInfo.GetEmail(),
			userInfo.GetAuthProvider(),
			&imageUrl,
			nil)
		err = a.userRepository.CreateUser(ctx, user)
		if err != nil {
			a.logger.Errorf(ctx, map[string]interface{}{
				"email": userInfo.GetEmail(),
			}, "Error create new user:", err)
		}
		a.logger.Info(ctx, map[string]interface{}{
			"email":   userInfo.GetEmail(),
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

func (a *AuthInteractorImpl) SignUp(ctx context.Context, input *presenter.SignUpInput) (*presenter.SignUpOuput, error) {
	//check if email exist
	existUser, err := a.userRepository.FindUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existUser != nil {
		return nil, exception.ErrEmailExist
	}
	newUser := models.NewUser(
		input.FirstName,
		input.LastName,
		input.Email,
		"basic",
		nil,
		&input.Password)
	//store user info to cache
	code, err := a.userConfirmService.StoreUserPending(ctx, newUser)
	if err != nil {
		a.logger.Errorf(ctx, map[string]interface{}{
			"email":   newUser.Email,
			"user_id": newUser.ID,
		}, "Error store user to cache: ", err)
		return nil, err
	}
	err = a.userConfirmService.SendMailConfirm(ctx, newUser, code)
	if err != nil {
		a.logger.Errorf(ctx, map[string]interface{}{
			"email":   newUser.Email,
			"user_id": newUser.ID,
		}, "Error send email to user: ", err)
		return nil, err
	}
	return &presenter.SignUpOuput{}, nil
}
