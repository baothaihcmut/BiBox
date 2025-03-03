package models

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                 primitive.ObjectID `bson:"_id"` // MongoDB automatically generates this ID if not provided
	FirstName          string             `bson:"first_name"`
	LastName           string             `bson:"last_name"`
	Email              string             `bson:"email"`
	Password           *string            `bson:"password,omitempty"`
	Image              *string            `bson:"image"`
	AuthProvider       string             `bson:"auth_provider"`
	CreatedAt          primitive.DateTime `bson:"created_at"`
	UpdatedAt          primitive.DateTime `bson:"updated_at"`
	CurrentStorageSize int                `bson:"current_storage_size"`
	LimitStorageSize   int                `bson:"limit_storage_size"`
}

// Constructor function to create a new user
func NewUser(firstName, lastName, email, authProvider string, image, password *string) *User {
	return &User{
		ID:                 primitive.NewObjectID(),
		FirstName:          firstName,
		LastName:           lastName,
		Email:              email,
		Password:           password,
		Image:              image,
		CreatedAt:          primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:          primitive.NewDateTimeFromTime(time.Now()),
		AuthProvider:       authProvider,
		CurrentStorageSize: 0,
		LimitStorageSize:   100 * (1024 ^ 3),
	}
}

func (u *User) IncreStorageSize(size int) error {
	newSize := u.CurrentStorageSize + size
	if newSize > u.LimitStorageSize {
		return exception.ErrStorageSizeExceedLimitSize
	}
	u.CurrentStorageSize = newSize
	return nil
}

func (u *User) DecreStorageSize(size int) error {
	if size > u.CurrentStorageSize {
		return exception.ErrStorageSizeLessThanZero
	}
	u.CurrentStorageSize = u.CurrentStorageSize - size
	return nil
}
