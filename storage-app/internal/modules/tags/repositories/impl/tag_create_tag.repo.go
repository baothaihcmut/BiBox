package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
)

func (m *MongoTagRepository) CreateTag(ctx context.Context, tag *models.Tag) error {
	_, err := m.collection.InsertOne(ctx, tag)
	if err != nil {
		m.logger.Errorf(ctx, map[string]any{
			"tag_id": tag.ID,
		}, "Error insert into tag collection: ", err)
		return err
	}
	return nil
}
