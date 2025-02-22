package interactors

import (
	"github.com/baothaihcmut/Storage-app/internal/modules/comment/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentInteractor struct {
	Repo *repositories.CommentRepository
}

// Pass a *mongo.Database instance as a parameter
func NewCommentInteractor(db *mongo.Database) *CommentInteractor {
	return &CommentInteractor{
		Repo: repositories.NewCommentRepository(db),
	}
}

func (ci *CommentInteractor) GetAllComments() ([]map[string]interface{}, error) {
	return ci.Repo.FetchComments()
}

func (ci *CommentInteractor) AddComment(fileID, userID, content string) error {
	return ci.Repo.CreateComment(fileID, userID, content)
}
