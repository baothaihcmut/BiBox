package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                 primitive.ObjectID `bson:"_id"` // MongoDB automatically generates this ID if not provided
	FirstName          string             `bson:"first_name"`
	LastName           string             `bson:"last_name"`
	Email              string             `bson:"email"`
	Image              string             `bson:"image"`
	CreatedAt          primitive.DateTime `bson:"created_at"`
	UpdatedAt          primitive.DateTime `bson:"updated_at"`
	CurrentStorageSize int                `bson:"current_storage_size"`
	LimitStorageSize   int                `bson:"limit_storage_size"`
}

// Constructor function to create a new user
func NewUser(firstName, lastName, email, image string) *User {
	return &User{
		ID:                 primitive.NewObjectID(),
		FirstName:          firstName,
		LastName:           lastName,
		Email:              email,
		Image:              image,
		CreatedAt:          primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:          primitive.NewDateTimeFromTime(time.Now()),
		CurrentStorageSize: 0,
		LimitStorageSize:   100,
	}
}
