package repositories

import (
	"context"

	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"github.com/baothaihcmut/Storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository interface {
	CreateFile(context.Context, *models.File) error
	FindFileById(ctx context.Context, id primitive.ObjectID, isDeleted bool) (*models.File, error)
}

type MongoFileRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func (f *MongoFileRepository) CreateFile(ctx context.Context, file *models.File) error {
	_, err := f.collection.InsertOne(ctx, file)
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error insert document:", err)
	}
	return nil
}

func (f *MongoFileRepository) FindFileById(ctx context.Context, id primitive.ObjectID, isDeleted bool) (*models.File, error) {
	var res models.File
	err := f.collection.FindOne(ctx, bson.M{
		"_id":        id,
		"is_deleted": isDeleted,
	}).Decode(&res)
	if err != nil {
		f.logger.Errorf(ctx, map[string]interface{}{
			"file_id":    id,
			"is_deleted": isDeleted,
		}, "Error find file by id:", err)
		return nil, err
	}
	return &res, nil
}
func NewMongoFileRepo(collection *mongo.Collection, logger logger.Logger) FileRepository {
	return &MongoFileRepository{
		collection: collection,
		logger:     logger,
	}
}
