package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/libs/pkg/consumer"
	"github.com/baothaihcmut/BiBox/libs/pkg/handler"
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/BiBox/libs/pkg/middlewares"
	"github.com/baothaihcmut/BiBox/libs/pkg/router"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/config"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/handlers"
	"github.com/sirupsen/logrus"

	"github.com/baothaihcmut/BiBox/storage-app-email/internal/services"
	"gopkg.in/gomail.v2"
)

type Server struct {
	consumer   *consumer.Consumer
	router     router.MessageRouter
	errHandler handler.ErrorHandler
	cfg        *config.CoreConfig
	mailDialer *gomail.Dialer
	logger     logger.Logger
}

func NewServer(
	mailDialer *gomail.Dialer,
	cfg *config.CoreConfig,
	logrus *logrus.Logger,
) *Server {
	errHandler := handler.NewErrorHandler(logrus)
	router := router.NewMessageRouter(errHandler)
	return &Server{
		router:     router,
		cfg:        cfg,
		mailDialer: mailDialer,
		errHandler: errHandler,
		consumer:   consumer.NewConsumer(router, &cfg.Consumer),
	}
}

func (s *Server) initApp() {
	//init global middleware
	s.router.RegisterGlobal(middlewares.ExtractHeaderMiddleware)
	s.router.RegisterGlobal(middlewares.LoggingMiddleware)
	//init service
	mailService := services.NewGmailService(s.mailDialer, &s.cfg.Mail)
	userMailService := services.NewUserMailService(mailService)
	//init handler
	userMailHandler := handlers.NewUserHandler(userMailService)
	userMailHandler.Init(s.router)
}
func (s *Server) Run() {
	s.initApp()
	consumerCfg := sarama.NewConfig()
	consumerCfg.Consumer.Return.Errors = true
	consumerCfg.Version = sarama.V3_3_1_0
	consumerGroup, err := sarama.NewConsumerGroup(s.cfg.Consumer.Brokers, s.cfg.Consumer.ConsumberGroupId, consumerCfg)
	if err != nil {
		s.logger.Error(context.Background(), nil, "Error running consumer")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		fmt.Println("Shutdown signal received. Closing consumer...")
		cancel()
	}()
	consumerRun := make(chan struct{}, 1)
	gracefullShutDown := make(chan struct{}, 1)
	go func() {
		s.router.Run(ctx, consumerRun)
		gracefullShutDown <- struct{}{}
	}()
	go func() {
		<-consumerRun
		for {
			if err := consumerGroup.Consume(ctx, s.cfg.Consumer.Topics, s.consumer); err != nil {
				fmt.Printf("Error consume message: %v\n", err)
			}
		}
	}()
	go func() {
		s.errHandler.Run(ctx)
	}()

	s.logger.Info(context.Background(), nil, "Consumer run...")
	<-ctx.Done()
	<-gracefullShutDown
	s.logger.Info(context.Background(), nil, "Consumer shutdown...")
}
