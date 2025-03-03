package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagOutput struct {
	Id   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}

func MaptoOuput(t *models.Tag) *TagOutput {
	return &TagOutput{
		Id:   t.ID,
		Name: t.Name,
	}
}
