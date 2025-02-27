package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/config"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/consumer"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/handlers"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/middlewares"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/router"
	"github.com/baothaihcmut/BiBox/storage-app-email/internal/services"
	"gopkg.in/gomail.v2"
)

type Server struct {
	consumer   *consumer.Consumer
	router     router.MessageRouter
	cfg        *config.CoreConfig
	mailDialer *gomail.Dialer
}

func NewServer(
	consumer *consumer.Consumer,
	router router.MessageRouter,
	mailDialer *gomail.Dialer,
	cfg *config.CoreConfig) *Server {
	return &Server{
		consumer:   consumer,
		router:     router,
		cfg:        cfg,
		mailDialer: mailDialer,
	}
}

func (s *Server) initApp() {
	//init global middleware
	s.router.RegisterGlobal(middlewares.ExtractEventMiddleware)
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
		fmt.Println("Error init consumer group")
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
	go func() {
		for {
			if err := consumerGroup.Consume(ctx, s.cfg.Consumer.Topics, s.consumer); err != nil {
				fmt.Printf("Error consume message: %v\n", err)
			}
			if ctx.Err() != nil {
				break
			}
		}
	}()
	fmt.Println("Consumer started")
	<-ctx.Done()
	fmt.Println("Shutting down consumer...")
	close(s.consumer.MsgChan)
	s.consumer.Wg.Wait()
	fmt.Println("Consumer shut down gracefully.")
}
