package impl

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
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
