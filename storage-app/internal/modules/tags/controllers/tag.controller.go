package controllers

import (
	"net/http"

	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
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

// @Sumary Find all tags
// @Description Find all tags
// @Sumary Find all tags
// @Description Find all tags
// @Tags tags
// @Accept json
// @Produce json
// @Param limit query int true "limit"
// @Param offset query int true "offset"
// @Param query query string true "query for search"
// @Success 201 {object} response.AppResponse{data=presenters.SearchTagsOutput} "Create file sucess, storage_detail.put_object_url is presign url for upload file"
// @Router   /tags/:id/files [get]
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

// @Sumary Find all file of tags
// @Description Find all file of tags
// @Sumary Find all file of tags
// @Description Find all file of tags
// @Tags tags
// @Accept json
// @Produce json
// @Param id path string true "id of tag must be object id"
// @Param limit query int true "limit"
// @Param offset query int true "offset"
// @Param sort_by query string true "sort"
// @Param is_asc query bool true "direction sort"
// @Success 201 {object} response.AppResponse{data=presenters.GetAllFileOfTagOutput} "Create file sucess, storage_detail.put_object_url is presign url for upload file"
// @Failure 403 {object} response.AppResponse{data=nil} "User don't have permission for this file operation"
// @Failure 404 {object} response.AppResponse{data=nil} "Tag of file not found"
// @Router   /tags/search [get]
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
