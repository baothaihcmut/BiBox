package controllers

import (
	"net/http"

	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	middleware "github.com/baothaihcmut/Bibox/storage-app/internal/common/middlewares"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/services"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/presenters"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type FileCommentController interface {
	Init(g *gin.RouterGroup)
}

type FileControllerImpl struct {
	interactor  interactors.FileCommentInteractor
	authHandler services.JwtService
	logger      logger.Logger
}

func (f *FileControllerImpl) handleCreateComment(c *gin.Context) {
	payload, _ := c.Get(string(constant.PayloadContext))
	res, err := f.interactor.CreateFileComment(c.Request.Context(), payload.(*presenters.CreateFileCommentInput))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.InitResponse(true, "Create file comment success", res))
}

func (f *FileControllerImpl) Init(g *gin.RouterGroup) {
	internal := g.Group("/files/:id/comments")
	internal.Use(middleware.AuthMiddleware(f.authHandler, f.logger, false))
	internal.POST("/add", middleware.ValidateMiddleware[presenters.CreateFileCommentInput](true, binding.JSON), f.handleCreateComment)
}

func NewFileCommentController(
	interactor interactors.FileCommentInteractor,
	authHandler services.JwtService,
	logger logger.Logger,
) *FileControllerImpl {
	return &FileControllerImpl{
		interactor:  interactor,
		authHandler: authHandler,
		logger:      logger,
	}
}
