package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (f *MongoFileRepository) FindFileById(ctx context.Context, id primitive.ObjectID, isDeleted bool) (*models.File, error) {
	var res models.File
	err := f.collection.FindOne(ctx, bson.M{
		"_id":        id,
		"is_deleted": isDeleted,
	}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		f.logger.Errorf(ctx, map[string]any{
			"file_id":    id,
			"is_deleted": isDeleted,
		}, "Error find file by id:", err)
		return nil, err
	}
	return &res, nil
}
