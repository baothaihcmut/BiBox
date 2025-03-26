package impl

import (
	"context"
	"errors"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoTagRepository) FindTagById(ctx context.Context, id primitive.ObjectID) (*models.Tag, error) {
	var res models.Tag
	err := m.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		m.logger.Errorf(ctx, map[string]any{
			"tag_id": id.Hex(),
		}, "Error find tag by id:", err)
		return nil, err
	}
	return &res, nil
}
