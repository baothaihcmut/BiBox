package controllers

import (
	"net/http"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	middleware "github.com/baothaihcmut/Bibox/storage-app/internal/common/middlewares"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/services"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/presenters"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TagController interface {
	Init(g *gin.RouterGroup)
}

type TagControllerImpl struct {
	interactor  interactors.TagInteractor
	authHandler services.JwtService
	logger      logger.Logger
}

func (f *TagControllerImpl) handleFindAllTags(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.GetAllTags(c.Request.Context(), payload.(*presenters.SerchTagsInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Find all tag success", res))
}
func (f *TagControllerImpl) handleGetAllFileOfTag(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.GetAllFileOfTag(c.Request.Context(), payload.(*presenters.GetAllFilOfTagInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Get all file of tag success", res))

}

func (f *TagControllerImpl) Init(g *gin.RouterGroup) {
	external := g.Group("/tags")
	external.GET("/search", middleware.ValidateMiddleware[presenters.SerchTagsInput](false, binding.Query), f.handleFindAllTags)
	internal := g.Group("/tags")
	internal.Use(middleware.AuthMiddleware(f.authHandler, f.logger, false))
	internal.GET("/:id/files", middleware.ValidateMiddleware[presenters.GetAllFilOfTagInput](true, binding.Query), f.handleGetAllFileOfTag)

}
func NewTagController(
	interactor interactors.TagInteractor,
	authHandler services.JwtService,
	logger logger.Logger,
) TagController {
	return &TagControllerImpl{
		interactor:  interactor,
		authHandler: authHandler,
		logger:      logger,
	}
}
