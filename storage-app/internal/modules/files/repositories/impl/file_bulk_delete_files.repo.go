package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *MongoFileRepository) BulkDeleteFile(ctx context.Context, files []*models.File) error {
	_, err := f.collection.DeleteMany(ctx, bson.M{
		"_id": bson.M{
			"$in": lo.Map(files, func(item *models.File, _ int) primitive.ObjectID {
				return item.ID
			}),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
