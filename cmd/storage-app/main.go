package main

import (
	"github.com/baothaihcmut/Storage-app/internal/config"
	"github.com/baothaihcmut/Storage-app/internal/server"
	"github.com/baothaihcmut/Storage-app/internal/server/initialize"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
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
	oauth2Google := initialize.InitializeOauth2(&config.Oauth2.Google, []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}, google.Endpoint)
	oauth2Github := initialize.InitializeOauth2(&config.Oauth2.Github, []string{
		"read:user", "user:email",
	}, github.Endpoint)

	//s3
	s3, err := initialize.InitalizeS3(config.S3)
	if err != nil {
		logger.Panic(err)
		panic(err)
	}
	s := server.NewServer(g, mongo, oauth2Google, oauth2Github, s3, logger, config)
	s.Run()
}
