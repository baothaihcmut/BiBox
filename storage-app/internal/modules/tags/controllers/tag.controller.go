package controllers

import (
	"net/http"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	middleware "github.com/baothaihcmut/Bibox/storage-app/internal/common/middlewares"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/presenters"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TagController interface {
	Init(g *gin.RouterGroup)
}

type TagControllerImpl struct {
	interactor interactors.TagInteractor
}

func (f *TagControllerImpl) handleFindAllTags(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.GetAllTags(c.Request.Context(), payload.(*presenters.SerchTagsInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Create file success", res))
}
func handleGetAllFileOfTag(c *gin.Context) {

}

func (f *TagControllerImpl) Init(g *gin.RouterGroup) {
	external := g.Group("/tags")
	external.GET("/search", middleware.ValidateMiddleware[presenters.SerchTagsInput](false, binding.Query), f.handleFindAllTags)
}
func NewTagController(interactor interactors.TagInteractor) TagController {
	return &TagControllerImpl{
		interactor: interactor,
	}
}
