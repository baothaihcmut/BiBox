package impl

import (
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTagRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewMongoTagRepository(collection *mongo.Collection, logger logger.Logger) *MongoTagRepository {
	return &MongoTagRepository{
		collection: collection,
		logger:     logger,
	}
}
