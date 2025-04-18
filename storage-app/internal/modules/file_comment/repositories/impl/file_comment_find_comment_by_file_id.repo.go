package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MongoFileCommentRepo) FindCommentByFileId(ctx context.Context, fileId primitive.ObjectID) ([]*models.FileComment, error) {
	cursor, err := m.collection.Find(ctx, bson.M{
		"file_id": fileId,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var res []*models.FileComment
	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}
