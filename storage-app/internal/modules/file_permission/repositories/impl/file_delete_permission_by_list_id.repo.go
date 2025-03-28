package impl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FilePermissionRepositoryImpl) DeletePermissionByListFileId(ctx context.Context, fileIds []primitive.ObjectID) error {
	_, err := f.collection.DeleteMany(ctx, bson.M{
		"file_id": bson.M{
			"$in": fileIds,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
