package controllers

import (
	"net/http"

	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	middleware "github.com/baothaihcmut/Bibox/storage-app/internal/common/middlewares"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/services"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/presenters"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController interface {
	Init(g *gin.RouterGroup)
}

type UserControllerImpl struct {
	userInteractor interactors.UserInteractor
	authHandler    services.JwtService
	logger         logger.Logger
}

func (u *UserControllerImpl) Init(g *gin.RouterGroup) {
	internal := g.Group("/users")
	internal.Use(middleware.AuthMiddleware(u.authHandler, u.logger, false))
	internal.GET("/search", middleware.ValidateMiddleware[presenters.SearchUserInput](false, binding.Query), u.handleSearchUserByEmail)

}

// @Sumary Search user by email
// @Description search user by email
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "email value can be any string"
// @Param offset query string false "off set must be greater than 0"
// @Param limit query string false "limit must be greater than 0"
// @Success 201 {object} response.AppResponse{data=presenters.SearchUserOuput} "search user success"
// @Router   /users/search [get]
func (u *UserControllerImpl) handleSearchUserByEmail(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := u.userInteractor.SearchUserByEmail(c.Request.Context(), payload.(*presenters.SearchUserInput))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Search user by email success", res))
}
func NewUserController(
	userInteractor interactors.UserInteractor,
	authHandler services.JwtService,
	logger logger.Logger,
) UserController {
	return &UserControllerImpl{
		userInteractor: userInteractor,
		authHandler:    authHandler,
		logger:         logger,
	}
}
