package repositories

import (
	"context"

	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"github.com/baothaihcmut/Storage-app/internal/modules/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(context.Context, *models.User) error
	FindUserByEmail(context.Context, string) (*models.User, error)
}

type MongoUserRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewMongoUserRepository(collection *mongo.Collection, logger logger.Logger) UserRepository {
	return &MongoUserRepository{collection: collection, logger: logger}
}
func (m *MongoUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := m.collection.InsertOne(ctx, user)
	if err != nil {
		m.logger.Errorf(ctx, map[string]interface{}{
			"component":  "repository",
			"user_email": user.Email,
		}, "Error insert document to Mongo:", err)
		return err
	}
	return nil
}

func (m *MongoUserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := m.collection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		m.logger.Errorf(ctx, map[string]interface{}{
			"component":  "repository",
			"user_email": user.Email,
		}, "Error find one document from Mongo:", err)
		return nil, err
	}
	return &user, err
}
