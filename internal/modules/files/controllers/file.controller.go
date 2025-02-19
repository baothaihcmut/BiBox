package controllers

import (
	"net/http"

	"github.com/baothaihcmut/Storage-app/internal/common/constant"
	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	middleware "github.com/baothaihcmut/Storage-app/internal/common/middlewares"
	"github.com/baothaihcmut/Storage-app/internal/common/response"
	"github.com/baothaihcmut/Storage-app/internal/modules/auth/services"
	"github.com/baothaihcmut/Storage-app/internal/modules/files/interactors"
	"github.com/baothaihcmut/Storage-app/internal/modules/files/presenters"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type FileController interface {
	Init(g *gin.RouterGroup)
}
type FileControllerImpl struct {
	interactor  interactors.FileInteractor
	authHandler services.JwtService
	logger      logger.Logger
}

func (f *FileControllerImpl) Init(g *gin.RouterGroup) {
	internal := g.Group("/files")
	internal.Use(middleware.AuthMiddleware(f.authHandler, f.logger, false))
	internal.POST("", middleware.ValidateMiddleware[presenters.CreateFileInput](false, binding.JSON), f.handleCreateFile)
	internal.PATCH("/uploaded/:id", middleware.ValidateMiddleware[presenters.UploadedFileInput](true), f.handleUploadedFile)

}

// @Sumary Create new file
// @Description Create new file
// @Tags files
// @Accept json
// @Produce json
// @Param file body presenters.CreateFileInput true "file information"
// @Success 201 {object} response.AppResponse{data=presenters.CreateFileOutput} "Create file sucess, storage_detail.put_object_url is presign url for upload file"
// @Failure 403 {object} response.AppResponse{data=nil} "User don't have permission for this file operation"
// @Failure 404 {object} response.AppResponse{data=nil} "Parent folder not found, Tag of file not found"
// @Router   /files [post]
func (f *FileControllerImpl) handleCreateFile(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.CreatFile(c.Request.Context(), payload.(*presenters.CreateFileInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Create file success", res))
}

func (f *FileControllerImpl) handleUploadedFile(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.UploadedFile(c.Request.Context(), payload.(*presenters.UploadedFileInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Uploaded file success", res))
}

func NewFileController(interactor interactors.FileInteractor, jwtService services.JwtService, logger logger.Logger) FileController {
	return &FileControllerImpl{
		interactor:  interactor,
		authHandler: jwtService,
		logger:      logger,
	}
}
