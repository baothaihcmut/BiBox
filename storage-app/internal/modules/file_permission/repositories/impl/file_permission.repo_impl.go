package impl

import (
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilePermissionRepositoryImpl struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewPermissionRepository(collection *mongo.Collection, logger logger.Logger) *FilePermissionRepositoryImpl {
	return &FilePermissionRepositoryImpl{
		collection: collection,
		logger:     logger,
	}
}
