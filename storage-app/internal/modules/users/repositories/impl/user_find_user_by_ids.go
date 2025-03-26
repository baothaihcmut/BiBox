package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MongoUserRepository) FindUserIdIds(ctx context.Context, ids []primitive.ObjectID) ([]*models.User, error) {
	cursor, err := m.collection.Find(ctx, bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []*models.User

	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}
