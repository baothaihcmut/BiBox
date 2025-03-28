package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (pr *FilePermissionRepositoryImpl) FindFilePermissionWithUser(ctx context.Context, fileId primitive.ObjectID) ([]*models.FilePermissionWithUser, error) {
	pipeline := mongo.Pipeline{
		//match stage
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "file_id", Value: fileId},
			}},
		},
		//lookup stage
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "localField", Value: "user_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "user"},
			}},
		},
		bson.D{{Key: "$unwind", Value: "$user"}},
		//projection stage
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "file_id", Value: 1},
				{Key: "user_id", Value: 1},
				{Key: "permission_type", Value: 1},
				{Key: "can_share", Value: 1},
				{Key: "access_secure_file", Value: 1},
				{Key: "user.first_name", Value: 1},
				{Key: "user.last_name", Value: 1},
				{Key: "user.email", Value: 1},
				{Key: "user.image", Value: 1},
			}},
		},
	}

	cursor, err := pr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []*models.FilePermissionWithUser
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil

}
