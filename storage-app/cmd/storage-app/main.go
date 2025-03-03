package main

import (
	"context"
	"log"

	"github.com/baothaihcmut/Bibox/storage-app/internal/config"

	"github.com/baothaihcmut/Bibox/storage-app/internal/server"
	"github.com/baothaihcmut/Bibox/storage-app/internal/server/initialize"

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
	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	// Initialize logger
	logger := initialize.InitializeLogger(&config.Logger)

	// Initialize Gin engine
	g := gin.Default()

	// Initialize MongoDB
	mongoClient, err := initialize.InitializeMongo(&config.Mongo)
	if err != nil {
		logger.Panic(err)
		log.Fatal("Failed to initialize MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())
	// Initialize OAuth2 (Google & GitHub)
	oauth2Google := initialize.InitializeOauth2(&config.Oauth2.Google, []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}, google.Endpoint)

	oauth2Github := initialize.InitializeOauth2(&config.Oauth2.Github, []string{
		"read:user", "user:email",
	}, github.Endpoint)

	// Initialize S3 Storage
	s3, err := initialize.InitalizeS3(config.S3)
	if err != nil {
		logger.Panic(err)
		log.Fatal("Failed to initialize S3:", err)
	}
	//kafka
	kafka, err := initialize.InitializeKafkaProducer(&config.Kafka)
	if err != nil {
		logger.Panic(err)
		panic(err)
	}
	defer kafka.Close()
	//redis
	redis, err := initialize.InitializeRedis(&config.Redis)
	if err != nil {
		logger.Panic(err)
		panic(err)
	}
	defer redis.Close()
	// Create a new server instance
	s := server.NewServer(g, mongoClient, oauth2Google, oauth2Github, s3, kafka, redis, logger, config)

	s.Run()
}
