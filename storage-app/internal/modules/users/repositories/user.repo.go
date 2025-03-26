package repositories

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(context.Context, *models.User) error
	UpdateUserStorageSize(context.Context, *models.User) error
	FindUserByEmail(context.Context, string) (*models.User, error)
	FindUserById(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	FindUserRegexAndCount(ctx context.Context, email string, limit, offset *int) ([]*models.User, int, error)
	FindUserIdIds(ctx context.Context, ids []primitive.ObjectID) ([]*models.User, error)
}
