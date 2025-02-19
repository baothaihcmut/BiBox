package main

import (
	"github.com/baothaihcmut/Storage-app/internal/config"
	"github.com/baothaihcmut/Storage-app/internal/server"
	"github.com/baothaihcmut/Storage-app/internal/server/initialize"
	"github.com/gin-gonic/gin"
)

// @title Storage App API
// @version 1.0
// @description This is a sample API for file storage
// @host localhost:8080
// @BasePath /api/v1
func main() {
	//config
	config, err := config.LoadConfig()
	if err != nil {
		panic("Error for load config")
	}

	//logger
	logger := initialize.InitializeLogger(&config.Logger)

	g := gin.Default()

	//mongo
	mongo, err := initialize.InitializeMongo(&config.Mongo)
	if err != nil {
		logger.Panic(err)
		panic(err)
	}

	//oauth 2 google
	oauth2 := initialize.InitializeOauth2(&config.Oauth2)

	//s3
	s3, err := initialize.InitalizeS3(config.S3)
	if err != nil {
		logger.Panic(err)
		panic(err)
	}
	s := server.NewServer(g, mongo, oauth2, s3, logger, config)
	s.Run()
}
