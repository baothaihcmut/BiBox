package interactors

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/mongo"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/presenter"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/services"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/repositories"
)

type AuthInteractor interface {
	ExchangeToken(context.Context, *presenter.ExchangeTokenInput) (*presenter.ExchangeTokenOutput, error)
	SignUp(context.Context, *presenter.SignUpInput) (*presenter.SignUpOutput, error)
	ConfirmSignUp(context.Context, *presenter.ConfirmSignUpInput) (*presenter.ConfirmSignUpOutput, error)
}

type AuthInteractorImpl struct {
	oauth2ServiceFactory services.Oauth2ServiceFactory
	jwtService           services.JwtService
	userRepository       repositories.UserRepository
	userConfirmService   services.UserConfirmService
	logger               logger.Logger
	mongService          mongo.MongoService
}

func NewAuthInteractor(
	oauth2 services.Oauth2ServiceFactory,
	userRepo repositories.UserRepository,
	jwtService services.JwtService,
	logger logger.Logger,
	userConfirmService services.UserConfirmService,
	mongoService mongo.MongoService,
) AuthInteractor {
	return &AuthInteractorImpl{
		userRepository:       userRepo,
		oauth2ServiceFactory: oauth2,
		jwtService:           jwtService,
		logger:               logger,
		userConfirmService:   userConfirmService,
		mongService:          mongoService,
	}
}

// ConfirmSignUp implements AuthInteractor.
func (a *AuthInteractorImpl) ConfirmSignUp(ctx context.Context, input *presenter.ConfirmSignUpInput) (*presenter.ConfirmSignUpOutput, error) {
	//find user by code
	user, err := a.userConfirmService.GetUserPedingConfirm(ctx, input.Code)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, exception.ErrInvalidConfirmCode
	}
	//store user to db
	session, err := a.mongService.BeginTransaction(ctx)
	if err != nil {
		a.logger.Errorf(ctx, map[string]any{
			"user_id":    user.ID,
			"user_email": user.Email,
		}, "Error start mongo session: ", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			a.mongService.RollbackTransaction(ctx, session)
		}
		a.mongService.EndTransansaction(ctx, session)
	}()
	err = a.userRepository.CreateUser(ctx, user)
	if err != nil {
		a.logger.Errorf(ctx, map[string]any{
			"user_id":    user.ID,
			"user_email": user.Email,
		}, "Error save user to db: ", err)
		return nil, err
	}
	//comit transaction
	err = a.mongService.CommitTransaction(ctx, session)
	if err != nil {
		return nil, err
	}
	//remove in cache
	err = a.userConfirmService.ConfirmSignUp(ctx, user, input.Code)
	if err != nil {
		return nil, err
	}
	return &presenter.ConfirmSignUpOutput{}, nil

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

func (a *AuthInteractorImpl) SignUp(ctx context.Context, input *presenter.SignUpInput) (*presenter.SignUpOutput, error) {
	//check user pending confirm
	isPendingConfirm, err := a.userConfirmService.IsUserPedingConfirm(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if isPendingConfirm {
		return nil, exception.ErrUserPedingSignUpConfirm
	}

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
	return &presenter.SignUpOutput{}, nil
}
