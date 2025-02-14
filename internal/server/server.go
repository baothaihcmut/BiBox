package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"github.com/baothaihcmut/Storage-app/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	g      *gin.Engine
	logger *logrus.Logger
	config *config.AppConfig
	mongo  *mongo.Client
}

func NewServer(g *gin.Engine, mongo *mongo.Client, logger *logrus.Logger, cfg *config.AppConfig) *Server {
	return &Server{
		g:      g,
		logger: logger,
		config: cfg,
		mongo:  mongo,
	}
}
func (s *Server) initApp() {
	//init cors
	s.g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,      // Allow cookies
		MaxAge:           24 * 3600, // Preflight cache duration
	}))

	logger := logger.NewLogger(s.logger)
	//global prefix
	globalGroup := s.g.Group("/api/v1")

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
