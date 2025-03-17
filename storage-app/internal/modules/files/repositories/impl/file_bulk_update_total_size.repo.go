package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (f *MongoFileRepository) BulkUpdateTotalSize(ctx context.Context, files []*models.File) error {
	updates := lo.Map(files, func(item *models.File, _ int) mongo.WriteModel {
		return mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": item.ID}).SetUpdate(bson.M{"$set": bson.M{"total_size": item.TotalSize}})
	})
	_, err := f.collection.BulkWrite(ctx, updates)
	if err != nil {
		return err
	}
	return nil
}
