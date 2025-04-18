package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoFileCommentRepo) FindCommentById(ctx context.Context, id primitive.ObjectID) (*models.FileComment, error) {
	var res models.FileComment
	if err := m.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}
