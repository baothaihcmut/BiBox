package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/samber/lo"
)

func (f *MongoFileRepository) BulkCreateFiles(ctx context.Context, files []*models.File) error {
	_, err := f.collection.InsertMany(ctx, lo.Map(files, func(item *models.File, _ int) interface{} {
		return item
	}))
	if err != nil {
		return err
	}
	return nil
}
