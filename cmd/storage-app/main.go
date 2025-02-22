package main

import (
	"github.com/baothaihcmut/Storage-app/internal/server"

	"github.com/baothaihcmut/Storage-app/internal/config"
	"github.com/baothaihcmut/Storage-app/internal/server/initialize"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	config, err := config.LoadConfig()
	if err != nil {
		panic("Error loading config")
	}

	// Initialize logger
	logger := initialize.InitializeLogger(&config.Logger)

	// Initialize Gin router
	g := gin.Default()

	// Initialize MongoDB
	mongo, err := initialize.InitializeMongo(&config.Mongo)
	if err != nil {
		logger.Panic(err)
		panic(err)
	}

	// Initialize OAuth2 for Google
	oauth2 := initialize.InitializeOauth2(&config.Oauth2)

	// Initialize server
	s := server.NewServer(g, mongo, oauth2, logger, config)

	// Run the server
	s.Run()
}
