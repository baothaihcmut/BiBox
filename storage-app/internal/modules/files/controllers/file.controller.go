package controllers

import (
	"net/http"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	middleware "github.com/baothaihcmut/Bibox/storage-app/internal/common/middlewares"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/services"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
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
	internal.POST("/add", middleware.ValidateMiddleware[presenters.CreateFileInput](false, binding.JSON), f.handleCreateFile)
	internal.PATCH("/:id/uploaded", middleware.ValidateMiddleware[presenters.UploadedFileInput](true), f.handleUploadedFile)
	internal.GET("", middleware.ValidateMiddleware[presenters.FindFileOfUserInput](false, binding.Query), f.handleFindFileOfUser)
	internal.GET("/:id/tags", middleware.ValidateMiddleware[presenters.GetFileTagsInput](true), f.handleGetTagOfFile)
	internal.GET("/:id/permissions", middleware.ValidateMiddleware[presenters.GetFilePermissionInput](true), f.handleGetPermissionOfFile)
	internal.GET("/:id/metadata", middleware.ValidateMiddleware[presenters.GetFileMetaDataInput](true), f.handleGetFileMetadata)
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
// @Router   /files/add [post]
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

// @Sumary Uploaded file
// @Description Uploaded file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "file id"
// @Success 201 {object} response.AppResponse{data=presenters.UploadedFileOutput} "Uploaded file sucess"
// @Failure 404 {object} response.AppResponse{data=nil} "file not found"
// @Failure 403 {object} response.AppResponse{data=nil} "file is folder"
// @Router   /files/:id/uploaded [patch]
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

// @Sumary Find file of user
// @Description Find file of user
// @Tags files
// @Accept json
// @Produce json
// @Param is_in_folder query bool false "file is in other folder, if null fetch all file"
// @Param is_folder query bool false "file is folder or not, if null fetch all file and folder"
// @Param sort_by query string true "sort field, allow short field: created_at, updated_at, opened_at"
// @Param is_asc query bool true "sort direction"
// @Param offset query int true "for pagination"
// @Param limit query int true "for pagination"
// @Param mime_type query string false "mime type of file, if is_folder is true not pass mime_type"
// @Success 200 {object} response.AppResponse{data=presenters.FindFileOfUserOuput} "Find file of user sucess"
// @Failure 400 {object} response.AppResponse{data=nil} "Un allow sort field, lack of query"
// @Router   /files [get]
func (f *FileControllerImpl) handleFindFileOfUser(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.FindAllFileOfUser(c.Request.Context(), payload.(*presenters.FindFileOfUserInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.InitResponse(true, "Find file of user success", res))
}

func (f *FileControllerImpl) handleGetTagOfFile(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.GetFileTags(c.Request.Context(), payload.(*presenters.GetFileTagsInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.InitResponse(true, "Find tags of file success", res))
}

func (f *FileControllerImpl) handleGetPermissionOfFile(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.GetFilePermissions(c.Request.Context(), payload.(*presenters.GetFilePermissionInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.InitResponse(true, "Find permissions of file success", res))
}

func (f *FileControllerImpl) handleGetFileMetadata(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.GetFileMetaData(c.Request.Context(), payload.(*presenters.GetFileMetaDataInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.InitResponse(true, "Find metadata of file success", res))
}

func NewFileController(interactor interactors.FileInteractor, jwtService services.JwtService, logger logger.Logger) FileController {
	return &FileControllerImpl{
		interactor:  interactor,
		authHandler: jwtService,
		logger:      logger,
	}
}
