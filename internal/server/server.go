package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	middleware "github.com/baothaihcmut/Storage-app/internal/common/middlewares"
	mongoLib "github.com/baothaihcmut/Storage-app/internal/common/mongo"
	"github.com/baothaihcmut/Storage-app/internal/common/storage"
	"github.com/baothaihcmut/Storage-app/internal/config"
	authController "github.com/baothaihcmut/Storage-app/internal/modules/auth/controllers"
	authInteractors "github.com/baothaihcmut/Storage-app/internal/modules/auth/interactors"
	authService "github.com/baothaihcmut/Storage-app/internal/modules/auth/services"
	fileController "github.com/baothaihcmut/Storage-app/internal/modules/files/controllers"
	fileInteractor "github.com/baothaihcmut/Storage-app/internal/modules/files/interactors"
	fileRepo "github.com/baothaihcmut/Storage-app/internal/modules/files/repositories"
	userRepo "github.com/baothaihcmut/Storage-app/internal/modules/users/repositories"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

type Server struct {
	g      *gin.Engine
	logger *logrus.Logger
	config *config.AppConfig
	mongo  *mongo.Client
	oauth2 *oauth2.Config
	s3     *s3.Client
}

func NewServer(g *gin.Engine, mongo *mongo.Client, oauth2 *oauth2.Config, s3 *s3.Client, logger *logrus.Logger, cfg *config.AppConfig) *Server {
	return &Server{
		g:      g,
		logger: logger,
		config: cfg,
		mongo:  mongo,
		oauth2: oauth2,
		s3:     s3,
	}
}
func (s *Server) initApp() {
	//init cors

	logger := logger.NewLogger(s.logger)

	//init repository
	userRepo := userRepo.NewMongoUserRepository(s.mongo.Database(s.config.Mongo.DatabaseName).Collection("users"), logger)
	fileRepo := fileRepo.NewMongoFileRepo(s.mongo.Database(s.config.Mongo.DatabaseName).Collection("files"), logger)
	//init service
	userJwtService := authService.NewUserJwtService(s.config.Jwt, logger)
	oauth2Service := authService.NewGoogleOauth2Service(s.oauth2, logger)
	storageService := storage.NewS3StorageService(s.s3, logger)
	mongoService := mongoLib.NewMongoTransactionService(s.mongo)

	//init interactor
	authInteractor := authInteractors.NewAuthInteractor(oauth2Service, userRepo, userJwtService, logger)
	fileInteractor := fileInteractor.NewFileInteractor(userRepo, fileRepo, logger, storageService, mongoService)
	//init controllers
	authController := authController.NewAuthController(authInteractor, &s.config.Jwt, &s.config.Oauth2)
	fileController := fileController.NewFileController(fileInteractor, userJwtService, logger)

	//global prefix
	globalGroup := s.g.Group("/api/v1")

	//init global middleware
	globalGroup.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,      // Allow cookies
		MaxAge:           24 * 3600, // Preflight cache duration
	}))
	globalGroup.Use(middleware.LoggingMiddleware(logger))
	globalGroup.Use(middleware.ErrorHandler())
	{
		authController.Init(globalGroup)
		fileController.Init(globalGroup)
	}
}

func (s *Server) Run() {
	s.initApp()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Server.Port))
	if err != nil {
		s.logger.Panic("Error init listener:", err)
	}
	go func() {
		if err := s.g.RunListener(listener); err != nil {
			s.logger.Panic("Error run gin engine:", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdown()
	<-ctx.Done()
}
