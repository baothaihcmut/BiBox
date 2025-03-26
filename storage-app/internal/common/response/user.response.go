package response

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserOutput struct {
	ID        primitive.ObjectID `json:"id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Image     *string            `json:"image"`
}

func MapToUserOutput(user *models.User) *UserOutput {
	return &UserOutput{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Image:     user.Image,
	}
}
