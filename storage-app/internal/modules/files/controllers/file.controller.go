package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	middleware "github.com/baothaihcmut/Bibox/storage-app/internal/common/middlewares"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	authService "github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/services"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type FileController interface {
	Init(g *gin.RouterGroup)
}
type FileControllerImpl struct {
	interactor               interactors.FileInteractor
	authHandler              authService.JwtService
	logger                   logger.Logger
	uploadProgressSSEManager services.FileUploadProgressSSEManagerService
}

func (f *FileControllerImpl) Init(g *gin.RouterGroup) {
	internal := g.Group("/files")
	internal.Use(middleware.AuthMiddleware(f.authHandler, f.logger, false))
	internal.POST("/add", middleware.ValidateMiddleware[presenters.CreateFileInput](false, binding.JSON), f.handleCreateFile)
	internal.POST("/upload-folder", middleware.ValidateMiddleware[presenters.UploadFolderInput](false, binding.JSON), f.handleUploadFolder)
	internal.PATCH("/:id/uploaded", middleware.ValidateMiddleware[presenters.UploadedFileInput](true), f.handleUploadedFile)
	internal.GET("/my-drive", middleware.ValidateMiddleware[presenters.GetAllFileOfUserInput](false, binding.Query), f.handleFindFileOfUser)
	internal.GET("/:id/tags", middleware.ValidateMiddleware[presenters.GetFileTagsInput](true), f.handleGetTagOfFile)
	internal.GET("/:id/permissions", middleware.ValidateMiddleware[presenters.GetFilePermissionInput](true), f.handleGetPermissionOfFile)
	internal.GET("/:id/metadata", middleware.ValidateMiddleware[presenters.GetFileMetaDataInput](true), f.handleGetFileMetadata)
	internal.GET("/:id/download-url", middleware.ValidateMiddleware[presenters.GetFileDownloadUrlInput](true, binding.Query), f.handleGetFileDownloadUrl)
	internal.GET("/:id/sub-file", middleware.ValidateMiddleware[presenters.GetSubFileOfFolderInput](true, binding.Query), f.handleGetSubFileOfFolder)
	internal.POST("/:id/permissions/add", middleware.ValidateMiddleware[presenters.AddFilePermissionInput](true, binding.JSON), f.handleAddFilePermission)
	internal.GET("/:id/my-permission", middleware.ValidateMiddleware[presenters.GetFilePermissionOfUserInput](true), f.handleGetFilePermissionOfUser)
	internal.PATCH("/:id/permissions/user/:userId/update", middleware.ValidateMiddleware[presenters.UpdateFilePermissionInput](true, binding.JSON), f.handleUpdateFilePermission)
	internal.DELETE("/:id/permissions/user/:userId/delete", middleware.ValidateMiddleware[presenters.DeleteFilePermissionOfUserInput](true), f.handleDeleteFilePermission)
	internal.PATCH("/:id/soft-delete", middleware.ValidateMiddleware[presenters.SoftDeleteFileInput](true), f.handleSoftDeleteFile)
	internal.PATCH("/:id/recover", middleware.ValidateMiddleware[presenters.RecoverFileInput](true, binding.JSON), f.handleRecoverFile)
	internal.DELETE("/:id/hard-delete", middleware.ValidateMiddleware[presenters.HardDeleteFileInput](true), f.handleHardDeleteFile)
	internal.GET("/:id/sse/upload-progress", f.handleSSEFileUploadProgress)
}

// @Sumary Create new file 1
// @Description Create new file 111
// @Sumary Create new file 1
// @Description Create new file 111
// @Tags files
// @Accept json
// @Produce json
// @Param file body presenters.CreateFileInput true "file information"
// @Success 201 {object} response.AppResponse{data=presenters.CreateFileOutput} "Create file sucess, storage_detail.put_object_url is presign url for upload file"
// @Failure 403 {object} response.AppResponse{data=nil} "User don't have permission for this file operation"
// @Failure 404 {object} response.AppResponse{data=nil} "Parent folder not found, Tag of file not found"
// @Router   /files/add [post]
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
// @Param id path string true "file id"
// @Success 201 {object} response.AppResponse{data=presenters.UploadedFileOutput} "Uploaded file sucess"
// @Failure 404 {object} response.AppResponse{data=nil} "file not found"
// @Failure 403 {object} response.AppResponse{data=nil} "file is folder"
// @Router   /files/:id/uploaded [patch]
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
// @Param is_deleted query bool false "if true file in bin, false is file in drive"
// @Param is_folder query bool false "file is folder or not, if null fetch all file and folder"
// @Param sort_by query string true "sort field, allow short field: created_at, updated_at, opened_at"
// @Param is_asc query bool true "sort direction"
// @Param offset query int true "for pagination"
// @Param limit query int true "for pagination"
// @Param mime_type query string false "mime type of file, if is_folder is true not pass mime_type"
// @Param mime_type query string false "mime type of file, if is_folder is true not pass mime_type"
// @Success 200 {object} response.AppResponse{data=presenters.GetAllFileOfUserOuput} "Find file of user sucess"
// @Failure 400 {object} response.AppResponse{data=nil} "Un allow sort field, lack of query"
// @Router   /files/my-drive [get]
// @Router   /files/my-drive [get]
func (f *FileControllerImpl) handleFindFileOfUser(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.GetAllFileOfUser(c.Request.Context(), payload.(*presenters.GetAllFileOfUserInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.InitResponse(true, "Find file of user success", res))
}

// @Sumary Get tag of file
// @Description Get tag of file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "file id"
// @Success 200 {object} response.AppResponse{data=presenters.GetFileTagsOutput} "Find tags of file sucess"
// @Failure 400 {object} response.AppResponse{data=nil} "miss id"
// @Failure 404 {object} response.AppResponse{data=nil} "file not found"
// @Failure 403 {object} response.AppResponse{data=nil} "permission denied"
// @Router   /files/:id/tags [get]
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

// @Sumary Get permission of file
// @Description Get permission of file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "file id"
// @Success 200 {object} response.AppResponse{data=presenters.GetFilePermissionOuput} "Find tags of file sucess"
// @Failure 400 {object} response.AppResponse{data=nil} "miss id"
// @Failure 404 {object} response.AppResponse{data=nil} "file not found"
// @Failure 403 {object} response.AppResponse{data=nil} "permission denied"
// @Router   /files/:id/permissions [get]
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

// @Sumary Get metadata of file
// @Description Get metadata of file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "file id"
// @Success 200 {object} response.AppResponse{data=presenters.GetFileMetaDataOuput} "Find tags of file sucess"
// @Failure 400 {object} response.AppResponse{data=nil} "miss id"
// @Failure 404 {object} response.AppResponse{data=nil} "file not found"
// @Failure 403 {object} response.AppResponse{data=nil} "permission denied"
// @Router   /files/:id/metadata [get]
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

// @Sumary Get download url of file
// @Description Get download url of file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "file id"
// @Param preview query string true "preview mode"
// @Success 200 {object} response.AppResponse{data=presenters.GetFileDownloadUrlOutput} "Find tags of file sucess"
// @Failure 400 {object} response.AppResponse{data=nil} "miss id"
// @Failure 404 {object} response.AppResponse{data=nil} "file not found"
// @Failure 409 {object} response.AppResponse{data=nil} "file is folder"
// @Failure 403 {object} response.AppResponse{data=nil} "permission denied"
// @Router   /files/:id/download-url [get]
func (f *FileControllerImpl) handleGetFileDownloadUrl(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.GetFileDownloadUrl(c.Request.Context(), payload.(*presenters.GetFileDownloadUrlInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.InitResponse(true, "Get file download url success", res))

}

// @Sumary Find file of user
// @Description Find file of user
// @Tags files
// @Accept json
// @Produce json
// @Param is_folder query bool false "file is folder or not, if null fetch all file and folder"
// @Param sort_by query string true "sort field, allow short field: created_at, updated_at, opened_at"
// @Param is_asc query bool true "sort direction"
// @Param is_deleted query bool true "true if file is deleted"
// @Param offset query int true "for pagination"
// @Param limit query int true "for pagination"
// @Param mime_type query string false "mime type of file, if is_folder is true not pass mime_type"
// @Param id path string false "file id"
// @Success 200 {object} response.AppResponse{data=presenters.GetSubFileOfFolderInput} "Find file of user sucess"
// @Failure 400 {object} response.AppResponse{data=nil} "Unallow sort field, lack of query"
// @Router   /files/:id/sub-file [get]
func (f *FileControllerImpl) handleGetSubFileOfFolder(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.GetAllSubFileOfFolder(c.Request.Context(), payload.(*presenters.GetSubFileOfFolderInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.InitResponse(true, "Get file structure success", res))
}

// @Sumary Upload folder
// @Description upload folder
// @Tags files
// @Accept json
// @Produce json
// @Param file body presenters.UploadFolderInput true "folder information"
// @Success 201 {object} response.AppResponse{data=presenters.UploadFolderOutput} "Create file sucess, storage_detail.put_object_url is presign url for upload file"
// @Failure 403 {object} response.AppResponse{data=nil} "User don't have permission for this file operation"
// @Failure 404 {object} response.AppResponse{data=nil} "Parent folder not found, Tag of file not found"
// @Router   /files/upload-folder [post]
func (f *FileControllerImpl) handleUploadFolder(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.UploadFolder(c.Request.Context(), payload.(*presenters.UploadFolderInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Upload folder success", res))
}

// @Sumary Add permission for file
// @Description Add permission for file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "file id must be UUID"
// @Param file body presenters.AddFilePermissionInput true "permission info"
// @Success 201 {object} response.AppResponse{data=presenters.AddFilePermissionOutput} "Add permission success"
// @Failure 403 {object} response.AppResponse{data=nil} "User don't have permission for this file operation"
// @Failure 404 {object} response.AppResponse{data=nil} "Parent folder not found, Tag of file not found"
// @Router   /files/:id/permissions/add [post]
func (f *FileControllerImpl) handleAddFilePermission(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.AddFilePermission(c.Request.Context(), payload.(*presenters.AddFilePermissionInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Add permission success", res))
}

// @Sumary Get file permission of user
// @Description Get file permission of user
// @Tags files
// @Accept json
// @Produce json
// @Param id path string false "file id"
// @Success 200 {object} response.AppResponse{data=presenters.GetSubFileOfFolderInput} "Find file of user sucess"
// @Failure 403 {object} response.AppResponse{data=nil} "User don't have permission for this file operation"
// @Failure 404 {object} response.AppResponse{data=nil} "Parent folder not found, Tag of file not found"
// @Router   /files/:id/my-permission [get]
func (f *FileControllerImpl) handleGetFilePermissionOfUser(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.GetFilePermissionOfUser(c.Request.Context(), payload.(*presenters.GetFilePermissionOfUserInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.InitResponse(true, "Get file permission of user success", res))
}

// @Sumary Update permission
// @Description Update permission
// @Tags files
// @Accept json
// @Produce json
// @Param id path string false "file id"
// @Param userId path string false "user id"
// @Param file body presenters.UpdateFilePermissionInput true "permission info"
// @Success 201 {object} response.AppResponse{data=presenters.UpdateFilePermissionOuput} "Update permission success"
// @Failure 400 {object} response.AppResponse{data=nil} "Unallow sort field, lack of query"
// @Router   /files/:id/permissions/user/:userId/update [patch]
func (f *FileControllerImpl) handleUpdateFilePermission(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.UpdateFilePermission(c.Request.Context(), payload.(*presenters.UpdateFilePermissionInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Update file permission of user success", res))
}

// @Sumary Delete permission
// @Description Delete permission
// @Tags files
// @Accept json
// @Produce json
// @Param id path string false "file id"
// @Param userId path string false "user id"
// @Success 204 {object} response.AppResponse{data=nil} "Delete permission success"
// @Failure 400 {object} response.AppResponse{data=nil} "Unallow sort field, lack of query"
// @Router   /files/:id/permissions/user/:userId/delete [delete]
func (f *FileControllerImpl) handleDeleteFilePermission(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.DeleteFilePermission(c.Request.Context(), payload.(*presenters.DeleteFilePermissionOfUserInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusNoContent, response.InitResponse(true, "Delete file permission of user success", res))
}

// @Sumary Soft delete file
// @Description Soft Delete file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string false "file id"
// @Success 201 {object} response.AppResponse{data=nil} "Delete permission success"
// @Failure 403 {object} response.AppResponse{data=nil} "Permission denied"
// @Router   /files/:id/soft-delete [patch]
func (f *FileControllerImpl) handleSoftDeleteFile(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.SoftDeleteFile(c.Request.Context(), payload.(*presenters.SoftDeleteFileInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Soft delete file  of user success", res))
}

// @Sumary Recover deleted file
// @Description Recover deleted file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string false "file id"
// @Success 201 {object} response.AppResponse{data=nil} "Delete permission success"
// @Failure 403 {object} response.AppResponse{data=nil} "Permission denied"
// @Router   /files/:id/recover [patch]
func (f *FileControllerImpl) handleRecoverFile(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.RecoverFile(c.Request.Context(), payload.(*presenters.RecoverFileInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Recover file  of user success", res))
}

// @Sumary Hard delete file
// @Description Hard Delete file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string false "file id"
// @Success 201 {object} response.AppResponse{data=nil} "Delete permission success"
// @Failure 403 {object} response.AppResponse{data=nil} "Permission denied"
// @Router   /files/:id/hard-delete [delete]
func (f *FileControllerImpl) handleHardDeleteFile(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.HardDeleteFile(c.Request.Context(), payload.(*presenters.HardDeleteFileInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Hard delete file  of user success", res))
}

func (f *FileControllerImpl) handleSSEFileUploadProgress(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	flusher, _ := c.Writer.(http.Flusher)
	sessionId := c.Query("session_id")
	msgCh, userId, err := f.uploadProgressSSEManager.Connect(c.Request.Context(), sessionId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	for msg := range msgCh {
		select {
		case <-c.Request.Context().Done():
			f.uploadProgressSSEManager.Disconnect(c.Request.Context(), userId, sessionId)
			return
		default:
			jsonData, err := json.Marshal(msg)
			if err != nil {
				c.Error(err)
				c.Abort()
				return
			}
			c.Writer.Write(jsonData)
			flusher.Flush()
		}
	}

}

func NewFileController(interactor interactors.FileInteractor,
	jwtService authService.JwtService,
	logger logger.Logger,
	sseManagerService services.FileUploadProgressSSEManagerService) FileController {
	return &FileControllerImpl{
		interactor:               interactor,
		authHandler:              jwtService,
		logger:                   logger,
		uploadProgressSSEManager: sseManagerService,
	}
}
