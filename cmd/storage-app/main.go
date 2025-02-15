package main

import (
	"github.com/baothaihcmut/Storage-app/internal/config"
	"github.com/baothaihcmut/Storage-app/internal/server"
	"github.com/baothaihcmut/Storage-app/internal/server/initialize"
	"github.com/gin-gonic/gin"
)

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

	s := server.NewServer(g, mongo, oauth2, logger, config)
	s.Run()
}
