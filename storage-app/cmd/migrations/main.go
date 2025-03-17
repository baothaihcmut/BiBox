package main

import (
	"context"
	"fmt"
	"log"

	"github.com/baothaihcmut/Bibox/storage-app/internal/config"
	"github.com/baothaihcmut/Bibox/storage-app/migrations"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var migrates = []func(context.Context, *mongo.Client, string) error{
	migrations.InsertTags,
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	mongoUri := cfg.Mongo.Uri
	database := cfg.Mongo.DatabaseName
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	ctx := context.Background()
	for _, migrate := range migrates {
		migrate(ctx, client, database)
	}
	fmt.Println("Done migrate111")
}
