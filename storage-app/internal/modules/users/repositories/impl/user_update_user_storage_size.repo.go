package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongoUserRepository) UpdateUserStorageSize(ctx context.Context, user *models.User) error {
	_, err := m.collection.UpdateOne(ctx, bson.M{
		"_id": user.ID,
	}, bson.M{
		"$set": bson.M{
			"current_storage_size": user.CurrentStorageSize,
		},
	})
	if err != nil {
		m.logger.Errorf(ctx, map[string]any{
			"storage_size": user.CurrentStorageSize,
			"user_id":      user.ID.Hex(),
		}, "Error update user storage size in mongo : ", err)
		return err
	}
	return nil
}
