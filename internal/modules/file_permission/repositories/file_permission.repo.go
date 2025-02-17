package repositories

import (
	"context"

	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"github.com/baothaihcmut/Storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilePermissionRepository interface {
	CreateFilePermission(context.Context, *models.FilePermission) (*models.FilePermission, error)
}

type MongoFilePermissionRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewMongoFilePermissionRepository(collection *mongo.Collection, logger logger.Logger) FilePermissionRepository {
	return &MongoFilePermissionRepository{
		collection: collection,
		logger:     logger,
	}
}
func (m *MongoFilePermissionRepository) CreateFilePermission(ctx context.Context, filePermission *models.FilePermission) (*models.FilePermission, error) {
	_, err := m.collection.InsertOne(ctx, filePermission)
	if err != nil {
		m.logger.Errorf(ctx, map[string]interface{}{
			"file_id": filePermission.FileID.Hex(),
			"user_id": filePermission.UserID.Hex(),
		}, "Error insert file permission: ", err)
		return nil, err
	}
	return filePermission, nil
}
