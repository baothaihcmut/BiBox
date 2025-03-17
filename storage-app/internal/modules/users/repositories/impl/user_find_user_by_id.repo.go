package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoUserRepository) FindUserById(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := m.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		m.logger.Errorf(ctx, map[string]any{
			"component":  "repository",
			"user_email": user.Email,
		}, "Error find user document from Mongo by id:", err)
		return nil, err
	}
	return &user, err
}
