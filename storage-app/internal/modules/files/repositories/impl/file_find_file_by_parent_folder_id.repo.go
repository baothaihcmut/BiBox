package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindFileByParentFolderId implements FileRepository.
func (f *MongoFileRepository) FindFileByParentFolderId(ctx context.Context, parentFolderId primitive.ObjectID) ([]*models.File, error) {
	filter := bson.M{
		"parent_folder_id": parentFolderId,
	}
	cursor, err := f.collection.Find(ctx, filter)
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error find file by folder id: %v", err)
		return nil, err
	}
	var res []*models.File
	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}
