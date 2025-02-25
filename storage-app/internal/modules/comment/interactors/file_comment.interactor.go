package interactors

// import (
// 	"storage-app/internal/modules/comment/repositories"
// )

// type CommentInteractor struct {
// 	Repo *repositories.CommentRepository
// }

// func NewCommentInteractor() *CommentInteractor {
// 	return &CommentInteractor{
// 		Repo: repositories.NewCommentRepository(),
// 	}
// }

// func (ci *CommentInteractor) GetAllComments() ([]map[string]interface{}, error) {
// 	return ci.Repo.FetchComments()
// }

// func (ci *CommentInteractor) AddComment(fileID, userID, content string) error {
// 	return ci.Repo.InsertComment(fileID, userID, content)
// }
