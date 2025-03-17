package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (f *MongoFileRepository) UpdateFile(ctx context.Context, file *models.File) error {

	_, err := f.collection.UpdateOne(ctx, bson.M{
		"_id": file.ID,
	}, bson.M{
		"$set": file,
	})
	if err != nil {
		f.logger.Errorf(ctx, map[string]any{
			"file_id": file.ID,
		}, "Error update file document:", err)
		return err
	}
	return nil
}
