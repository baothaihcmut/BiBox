package impl

import (
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewMongoUserRepository(collection *mongo.Collection, logger logger.Logger) repositories.UserRepository {
	return &MongoUserRepository{collection: collection, logger: logger}
}
