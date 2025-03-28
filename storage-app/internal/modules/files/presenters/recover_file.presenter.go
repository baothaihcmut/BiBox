package presenters

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecoverFileInput struct {
	Id                  string              `uri:"id" validate:"required"`
	DestinationFolderId *primitive.ObjectID `json:"dest_folder_id"`
}

type RecoverFileOutput struct {
	Files []*response.FileOutput `json:"files"`
}
