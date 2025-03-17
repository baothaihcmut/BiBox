package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
)

func (f *MongoFileRepository) CreateFile(ctx context.Context, file *models.File) error {
	_, err := f.collection.InsertOne(ctx, file)
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error insert document:", err)
	}
	return nil
}
