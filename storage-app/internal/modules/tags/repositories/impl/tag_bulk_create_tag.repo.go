package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
)

func (m *MongoTagRepository) BulkCreateTag(ctx context.Context, tags []*models.Tag) error {
	input := make([]interface{}, len(tags))
	for idx, tag := range tags {
		input[idx] = tag
	}
	_, err := m.collection.InsertMany(ctx, input)
	if err != nil {
		m.logger.Errorf(ctx, nil, "Error insert many into tag collection: ", err)
		return err
	}
	return nil
}
