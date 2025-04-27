package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/baothaihcmut/BiBox/libs/pkg/consumer"
	"github.com/baothaihcmut/BiBox/libs/pkg/handler"
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/BiBox/libs/pkg/router"
	"github.com/baothaihcmut/Bibox/storage-app/docs"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/cache"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/lock"
	middleware "github.com/baothaihcmut/Bibox/storage-app/internal/common/middlewares"
	mongoLib "github.com/baothaihcmut/Bibox/storage-app/internal/common/mongo"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/monitor"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/queue"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/storage"
	"github.com/baothaihcmut/Bibox/storage-app/internal/config"
	authController "github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/controllers"
	authInteractors "github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/interactors"
	authService "github.com/baothaihcmut/Bibox/storage-app/internal/modules/auth/services"
	commentController "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/controllers"
	commentInteractor "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/interactors/impl"
	commentRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/repositories/impl"
	filePermssionRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories/impl"
	filePermissionService "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/services"
	fileController "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/controllers"
	fileInteractor "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/interactors/impl"
	fileRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/repositories/impl"
	fileService "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services"
	fileServie "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services/impl"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/handlers"
	notificationRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/repositories/impl"
	notificationSvc "github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/services/impl"
	tagController "github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/controllers"
	tagInteractor "github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/interactors/impl"
	tagRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/repositories/impl"

	userController "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/controllers"
	userInteractor "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/interactors/impl"
	userRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/repositories/impl"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

type Server struct {
	consumerErrHandler handler.ErrorHandler
	msgRouter          router.MessageRouter
	g                  *gin.Engine
	logger             *logrus.Logger
	config             *config.AppConfig
	mongo              *mongo.Client
	googleOauth2       *oauth2.Config
	githubOauth2       *oauth2.Config
	s3                 *s3.Client
	kafkaProducer      sarama.SyncProducer
	redis              *redis.Client
}

func NewServer(
	g *gin.Engine,
	mongo *mongo.Client,
	googleoauth2 *oauth2.Config,
	githubOauth2 *oauth2.Config,
	s3 *s3.Client,
	kafkaProducer sarama.SyncProducer,
	redis *redis.Client,
	logger *logrus.Logger,
	cfg *config.AppConfig) *Server {
	errHandler := handler.NewErrorHandler(logger)
	msgRouter := router.NewMessageRouter(errHandler)
	return &Server{
		consumerErrHandler: errHandler,
		msgRouter:          msgRouter,
		g:                  g,
		logger:             logger,
		config:             cfg,
		mongo:              mongo,
		googleOauth2:       googleoauth2,
		githubOauth2:       githubOauth2,
		kafkaProducer:      kafkaProducer,
		redis:              redis,
		s3:                 s3,
	}
}
func (s *Server) initApp() {
	//init cors

	logger := logger.NewLogger(s.logger)
	//init external service
	kafkaService := queue.NewKafkaService(s.kafkaProducer)
	redisService := cache.NewRedisService(s.redis)
	userConfirmService := authService.NewUserConfirmService(redisService, kafkaService, logger)
	//init repository
	userRepo := userRepo.NewMongoUserRepository(s.mongo.Database(s.config.Mongo.DatabaseName).Collection("users"), logger)
	fileRepo := fileRepo.NewMongoFileRepo(s.mongo.Database(s.config.Mongo.DatabaseName).Collection("files"), logger)
	tagRepo := tagRepo.NewMongoTagRepository(s.mongo.Database(s.config.Mongo.DatabaseName).Collection("tags"), logger)
	filePermssionRepo := filePermssionRepo.NewPermissionRepository(s.mongo.Database(s.config.Mongo.DatabaseName).Collection("file_permissions"), logger)
	fileCommentRepo := commentRepo.NewMongoFileCommentRepo(s.mongo.Database(s.config.Mongo.DatabaseName).Collection("file_comments"))
	notificationRepo := notificationRepo.NewNotificationRepo(s.mongo.Database(s.config.Mongo.DatabaseName).Collection("notifications"))
	//init service
	userJwtService := authService.NewUserJwtService(s.config.Jwt, logger)
	googleOauth2Service := authService.NewGoogleOauth2Service(s.googleOauth2, logger)
	githubOauth2Service := authService.NewGithubOauth2Service(s.githubOauth2, logger)
	oauth2SerivceFactory := authService.NewOauth2ServiceFactory()
	oauth2SerivceFactory.Register(authService.GoogleOauth2Token, googleOauth2Service)
	oauth2SerivceFactory.Register(authService.GithubOauth2Token, githubOauth2Service)
	storageService := storage.NewS3StorageService(s.s3, logger, &s.config.S3)
	mongoService := mongoLib.NewMongoTransactionService(s.mongo)
	passwordService := authService.NewPasswordService()
	filePermssionService := filePermissionService.NewPermissionService(filePermssionRepo)
	fileStructureService := fileService.NewFileStructureService()
	notificationService := notificationSvc.NewNotificationService(notificationRepo, kafkaService, logger)
	fileProgressService := fileServie.NewFileUploadProgressService(redisService, logger)
	notificationSSEManagerService := notificationSvc.NewNotificationSSEManagerService(
		redisService,
		logger,
	)
	fileUploadProgressSSEManager := fileServie.NewFileUploadProgressSSEManagerService(fileProgressService, redisService, logger)
	distributedLockService := lock.NewRedisDistributedLockService(s.redis)
	//init interactor
	userInteractor := userInteractor.NewUserInteractor(userRepo)
	authInteractor := authInteractors.NewAuthInteractor(oauth2SerivceFactory, userRepo, userJwtService, logger, userConfirmService, mongoService, passwordService)
	fileInteractor := fileInteractor.NewFileInteractor(userRepo, tagRepo, fileRepo, filePermssionService, filePermssionRepo, fileStructureService, notificationService, fileProgressService, logger, storageService, mongoService, distributedLockService)
	tagInteractor := tagInteractor.NewTagInteractor(tagRepo, fileRepo, logger, mongoService)
	fileCommentInteractor := commentInteractor.NewFileCommentInteractor(fileCommentRepo, fileRepo, userRepo, filePermssionService, mongoService, logger)

	//init controllers
	authController := authController.NewAuthController(authInteractor, &s.config.Jwt, &s.config.Oauth2)
	fileController := fileController.NewFileController(fileInteractor, userJwtService, logger, fileUploadProgressSSEManager)
	userController := userController.NewUserController(userInteractor, userJwtService, logger)
	tagController := tagController.NewTagController(tagInteractor, userJwtService, logger)
	fileCommentController := commentController.NewFileCommentController(fileCommentInteractor, userJwtService, logger)

	//init event handler
	notifictionEventHandler := handlers.NewNotificationEventHandler(
		notificationSSEManagerService,
	)

	//register message router
	notifictionEventHandler.Init(s.msgRouter)

	//register metrics monitor
	httpRequestTotalMetric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total http request",
		},
		[]string{"method", "uri", "status"},
	)
	httpRequestDurationMetric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_requests_duration",
			Help:    "Duration http request",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "uri", "status"},
	)
	prometheusService := monitor.NewPrometheusService(
		httpRequestTotalMetric,
		httpRequestDurationMetric,
	)

	//init global middleware
	s.g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://storage-app-web.spsohcmut.xyz", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * 3600,
	}))

	// Global middleware
	s.g.Use(middleware.LoggingMiddleware(logger))
	s.g.Use(middleware.PrometheuseMiddleware(prometheusService))
	s.g.Use(middleware.ErrorHandler())
	s.g.GET("/metrics", gin.WrapH(promhttp.Handler()))

	//global prefix
	globalGroup := s.g.Group("/api/v1")

	//swagger
	globalGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/v1"
	{
		authController.Init(globalGroup)
		fileController.Init(globalGroup)
		userController.Init(globalGroup)
		tagController.Init(globalGroup)
		fileCommentController.Init(globalGroup)
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
	//for consumer

	consumer := consumer.NewConsumer(s.msgRouter, &s.config.Consumer)
	consumerCfg := sarama.NewConfig()
	consumerCfg.Consumer.Return.Errors = true
	consumerCfg.Version = sarama.V2_7_0_0
	consumerGroup, err := sarama.NewConsumerGroup(s.config.Consumer.Brokers, s.config.Consumer.ConsumberGroupId, consumerCfg)
	if err != nil {
		s.logger.Error(context.Background(), nil, "Error running consumer")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumerRun := make(chan struct{}, 1)
	gracefullShutDown := make(chan struct{}, 1)
	go func() {
		s.msgRouter.Run(ctx, consumerRun)
		gracefullShutDown <- struct{}{}
	}()
	go func() {
		<-consumerRun
		for {
			if err := consumerGroup.Consume(ctx, s.config.Consumer.Topics, consumer); err != nil {
				fmt.Printf("Error consume message: %v\n", err)
			}
		}
	}()
	go func() {
		s.consumerErrHandler.Run(ctx)
	}()

	s.logger.Info("Consumer run...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	//cancel consumer
	cancel()

	ctx, shutdown := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdown()
	<-ctx.Done()
}
