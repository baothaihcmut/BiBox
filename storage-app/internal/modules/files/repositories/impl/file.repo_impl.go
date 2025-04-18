package impl

import (
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoFileRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewMongoFileRepo(collection *mongo.Collection, logger logger.Logger) repositories.FileRepository {
	return &MongoFileRepository{
		collection: collection,
		logger:     logger,
	}
}
