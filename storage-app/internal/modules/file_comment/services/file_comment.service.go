package services

import (
	"context"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	fileObjectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return &InvalidInputError{"Invalid file ID format"}
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return &InvalidInputError{"Invalid user ID format"}
	}

	// context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// call repository
	return cs.Repo.CreateComment(ctx, fileObjectID, userObjectID, content)
}

func (cs *CommentService) AnswerComment(commentID, userID, content string) error {
	if commentID == "" || userID == "" || content == "" {
		return &InvalidInputError{"All fields are required"}
	}

	commentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return &InvalidInputError{"Invalid comment ID format"}
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return &InvalidInputError{"Invalid user ID format"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call repository
	return cs.Repo.AnswerComment(ctx, commentObjectID, userObjectID, content)
}

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

type PermissionError struct {
	Message string
}

func (e *PermissionError) Error() string {
	return e.Message
}
