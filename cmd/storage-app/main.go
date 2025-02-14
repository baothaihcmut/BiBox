package main

import (
	"io"

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

	//gin
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	g := gin.Default()

	//mongo
	mongo, err := initialize.InitializeMongo(&config.Mongo)
	if err != nil {
		logger.Panic(err)
		panic(err)
	}

	s := server.NewServer(g, mongo, logger, config)
	s.Run()
}
