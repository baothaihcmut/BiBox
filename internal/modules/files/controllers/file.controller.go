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
	internal.POST("", f.handleCreateFile)

}
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

func NewFileController(interactor interactors.FileInteractor, jwtService services.JwtService, logger logger.Logger) FileController {
	return &FileControllerImpl{
		interactor:  interactor,
		authHandler: jwtService,
		logger:      logger,
	}
}
