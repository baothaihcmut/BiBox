package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
)

func (m *MongoUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := m.collection.InsertOne(ctx, user)
	if err != nil {
		m.logger.Errorf(ctx, map[string]any{
			"component":  "repository",
			"user_email": user.Email,
		}, "Error insert document to Mongo:", err)
		return err
	}
	return nil
}
