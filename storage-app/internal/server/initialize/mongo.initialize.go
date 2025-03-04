package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitializeMongo(cfg *config.MongoConfig) (*mongo.Client, error) {
	clientOption := options.Client().
		ApplyURI(cfg.Uri).
		SetMaxPoolSize(uint64(cfg.MaxPoolSize)).
		SetMinPoolSize(uint64(cfg.MinPoolSize)).
		SetConnectTimeout(time.Second * time.Duration(cfg.ConnectionTimeout))
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		return nil, fmt.Errorf("error connect to MongoDb: %v", err)
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %v", err)
	}
	return client, nil
}
