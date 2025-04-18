package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (u *MongoUserRepository) FindUsersByEmailList(ctx context.Context, emails []string) ([]*models.User, error) {
	cursor, err := u.collection.Find(ctx, bson.M{
		"email": bson.M{
			"$in": emails,
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
