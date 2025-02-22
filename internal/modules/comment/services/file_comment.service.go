package services

import (
	"storage-app/internal/modules/comment/repositories"
)

type CommentService struct {
	Repo *repositories.CommentRepository
}

func NewCommentService() *CommentService {
	return &CommentService{
		Repo: repositories.NewCommentRepository(),
	}
}

func (cs *CommentService) GetAllComments() ([]map[string]interface{}, error) {
	return cs.Repo.FetchComments()
}

func (cs *CommentService) AddComment(fileID, userID, content string) error {
	// Business logic, validation, etc.
	if fileID == "" || userID == "" || content == "" {
		return &InvalidInputError{"All fields are required"}
	}
	return cs.Repo.InsertComment(fileID, userID, content)
}

// Custom error handling
type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}
