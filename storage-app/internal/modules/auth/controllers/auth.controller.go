package controllers

import (
	"fmt"
	"net/http"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	middleware "github.com/baothaihcmut/Bibox/storage-app/internal/common/middlewares"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/config"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/presenter"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController interface {
	Init(g *gin.RouterGroup)
}

type AuthControllerImpl struct {
	authInteractor interactors.AuthInteractor
	jwtConfig      *config.JwtConfig
	oauth2Config   *config.Oauth2Config
}

// @Sumary Exchange Google token
// @Description Exchange Google auth code
// @Tags auth
// @Accept json
// @Produce json
// @Param authCode body presenter.ExchangeTokenInput true "auth code from google oauth2 resposne"
// @Success 201 {object} response.AppResponse{data=nil} "Login success"
// @Failure 401 {object} response.AppResponse{data=nil} "Wrong auth code"
//
//	@Router   /auth/exchange [post]
func (a *AuthControllerImpl) handleExchangeToken(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := a.authInteractor.ExchangeToken(c.Request.Context(), payload.(*presenter.ExchangeTokenInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	//set cookie
	c.SetCookie("access_token", res.AccessToken, a.jwtConfig.AccessToken.Age*60*60, "/", " spsohcmut.xyz", true, true)
	c.SetCookie("refresh_token", res.RefreshToken, a.jwtConfig.RefreshToken.Age*60*60, "/", "spsohcmut.xyz", true, true)
	c.JSON(http.StatusCreated, response.InitResponse[any](true, "Login sucess", nil))
}

// @Sumary Sign up
// @Description Sign up
// @Tags auth
// @Accept json
// @Produce json
// @Param request body presenter.SignUpInput true "information for sign up"
// @Success 201 {object} response.AppResponse{data=presenter.SignUpOutput} "Sign up success"
// @Failure 409 {object} response.AppResponse{data=nil} "Email exist, email is pending for cofirm"
//
//	@Router   /auth/sign-up [post]
func (a *AuthControllerImpl) handleSignUp(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := a.authInteractor.SignUp(c.Request.Context(), payload.(*presenter.SignUpInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Sign up sucess, please check your email for confirmation", res))
}

// @Sumary Confirm sign up
// @Description Confirm sign up
// @Tags auth
// @Accept json
// @Produce json
// @Param request body presenter.ConfirmSignUpInput true "code for confirm"
// @Success 201 {object} response.AppResponse{data=presenter.ConfirmSignUpOutput} "Confirm sign up success"
// @Failure 401 {object} response.AppResponse{data=nil} "Invalid confirm code"
//
//	@Router   /auth/confirm [post]
func (a *AuthControllerImpl) handleConfirmSignUp(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := a.authInteractor.ConfirmSignUp(c.Request.Context(), payload.(*presenter.ConfirmSignUpInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Confirm sign up success, please login again", res))
}

// @Sumary Log in
// @Description Log in
// @Tags auth
// @Accept json
// @Produce json
// @Param request body presenter.LogInInput true "information for log in"
// @Success 201 {object} response.AppResponse{data=nil} "Login success"
// @Failure 401 {object} response.AppResponse{data=nil} "Wrong password or email"
//
//	@Router   /auth/log-in [post]
func (a *AuthControllerImpl) handleLogIn(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := a.authInteractor.LogIn(c.Request.Context(), payload.(*presenter.LogInInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.SetCookie("access_token", res.AccessToken, a.jwtConfig.AccessToken.Age*60*60, "/", " spsohcmut.xyz", true, true)
	c.SetCookie("refresh_token", res.RefreshToken, a.jwtConfig.RefreshToken.Age*60*60, "/", "spsohcmut.xyz", true, true)
	c.JSON(http.StatusCreated, response.InitResponse[any](true, "Login sucess", nil))
}

func (a *AuthControllerImpl) Init(g *gin.RouterGroup) {
	external := g.Group("/auth")
	external.POST(
		"/exchange",
		middleware.ValidateMiddleware[presenter.ExchangeTokenInput](false, binding.JSON),
		a.handleExchangeToken,
	)
	external.POST(
		"/sign-up",
		middleware.ValidateMiddleware[presenter.SignUpInput](false, binding.JSON),
		a.handleSignUp,
	)
	external.POST(
		"/confirm",
		middleware.ValidateMiddleware[presenter.ConfirmSignUpInput](false, binding.JSON),
		a.handleConfirmSignUp,
	)
	external.POST(
		"/log-in",
		middleware.ValidateMiddleware[presenter.LogInInput](false, binding.JSON),
		a.handleLogIn,
	)
	external.GET("/callback", func(c *gin.Context) {
		authCode := c.Query("code")
		if authCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing auth code"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": authCode})
	})

	external.GET("/redirect", func(c *gin.Context) {
		provider := c.Query("provider")
		var authRedirectURL string
		switch provider {
		case "google":
			authRedirectURL = fmt.Sprintf(
				"https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&access_type=offline&prompt=consent",
				a.oauth2Config.Google.ClientId, a.oauth2Config.Google.RedirectURI, "email profile",
			)
		case "github":
			authRedirectURL = fmt.Sprintf(
				"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user user:email",
				a.oauth2Config.Github.ClientId, a.oauth2Config.Github.RedirectURI)
		}

		// Redirect user to Google OAuth
		c.Redirect(http.StatusFound, authRedirectURL)
	})

}

func NewAuthController(interactor interactors.AuthInteractor, jwtConfig *config.JwtConfig, oauth2Config *config.Oauth2Config) AuthController {
	return &AuthControllerImpl{
		authInteractor: interactor,
		jwtConfig:      jwtConfig,
		oauth2Config:   oauth2Config,
	}
}
