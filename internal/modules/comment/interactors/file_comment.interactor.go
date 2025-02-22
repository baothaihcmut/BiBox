package interactors

import (
	"github.com/baothaihcmut/Storage-app/internal/modules/comment/repositories"
)

type CommentInteractor struct {
	Repo *repositories.CommentRepository
}

func NewCommentInteractor(repo *repositories.CommentRepository) *CommentInteractor {
	return &CommentInteractor{
		Repo: repo,
	}
}
func (ci *CommentInteractor) GetAllComments() ([]map[string]interface{}, error) {
	return ci.Repo.FetchComments()
}

func (ci *CommentInteractor) AddComment(fileID, userID, content string) error {
	return ci.Repo.CreateComment(fileID, userID, content)
}
