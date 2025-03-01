package services

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentService struct {
	Repo *repositories.CommentRepository
}

func NewCommentService(db *mongo.Database) *CommentService {
	return &CommentService{
		Repo: repositories.NewCommentRepository(db),
	}
}

func (cs *CommentService) GetAllComments() ([]map[string]interface{}, error) {
	return cs.Repo.FetchComments()
}

func (cs *CommentService) AddComment(fileID, userID, content string) error {
	if fileID == "" || userID == "" || content == "" {
		return &InvalidInputError{"All fields are required"}
	}

	return cs.Repo.CreateComment(fileID, userID, content)
}

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}
