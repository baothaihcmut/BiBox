package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MongoTagRepository) FindAllTagInList(ctx context.Context, tagIds []primitive.ObjectID) ([]*models.Tag, error) {
	cursor, err := m.collection.Find(ctx, bson.D{{
		Key: "_id", Value: bson.D{{
			Key:   "$in",
			Value: tagIds,
		}},
	}})
	if err != nil {
		m.logger.Errorf(ctx, map[string]any{
			"tag_ids": tagIds,
		}, "Error find all tag in list: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var tags []*models.Tag
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, err
	}
	return tags, err

}
