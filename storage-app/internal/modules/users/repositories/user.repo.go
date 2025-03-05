package repositories

import (
	"context"
	"fmt"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(context.Context, *models.User) error
	UpdateUserStorageSize(context.Context, *models.User) error
	FindUserByEmail(context.Context, string) (*models.User, error)
	FindUserById(ctx context.Context, id primitive.ObjectID) (*models.User, error)
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
		m.logger.Errorf(ctx, map[string]any{
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
		m.logger.Errorf(ctx, map[string]any{
			"component":  "repository",
			"user_email": user.Email,
		}, "Error find user document from Mongo by email:", err)
		return nil, err
	}
	return &user, err
}

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
func (m *MongoUserRepository) FindUserById(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := m.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&user)
	fmt.Println(err)
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
