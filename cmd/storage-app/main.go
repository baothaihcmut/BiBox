package main

import (
	"log"

	"github.com/baothaihcmut/Storage-app/internal/config"
	commentControllers "github.com/baothaihcmut/Storage-app/internal/modules/comment/controllers"
	commentInteractors "github.com/baothaihcmut/Storage-app/internal/modules/comment/interactors"
	commentRepo "github.com/baothaihcmut/Storage-app/internal/modules/comment/repositories"
	permControllers "github.com/baothaihcmut/Storage-app/internal/modules/permission/controllers"
	permInteractors "github.com/baothaihcmut/Storage-app/internal/modules/permission/interactors"
	"github.com/baothaihcmut/Storage-app/internal/server"
	"github.com/baothaihcmut/Storage-app/internal/server/initialize"

	"github.com/gin-gonic/gin"
)

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

	// Select the database
	mongoDatabase := mongoClient.Database(config.Mongo.DatabaseName)

	// Initialize OAuth2 for Google authentication
	oauth2 := initialize.InitializeOauth2(&config.Oauth2)

	// Create a new server instance
	s := server.NewServer(g, mongoClient, oauth2, logger, config)

	permissionInteractor := permInteractors.NewPermissionInteractor(mongoDatabase)
	permissionController := permControllers.NewPermissionController(permissionInteractor)

	commentRepository := commentRepo.NewCommentRepository(mongoDatabase)
	commentInteractor := commentInteractors.NewCommentInteractor(commentRepository)
	commentController := commentControllers.NewCommentController(commentInteractor)

	server.SetupRoutes(g, permissionController, commentController)

	go s.Run()

}
