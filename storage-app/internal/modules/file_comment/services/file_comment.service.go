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

func (cs *CommentService) GetAllComments() ([]map[string]any, error) {
	return cs.Repo.FetchComments()
}

func (cs *CommentService) AddComment(fileID, userID, content string) error {
	if fileID == "" || userID == "" || content == "" {
		return &InvalidInputError{"All fields are required"}
	}

	// Convert fileID and userID to primitive.ObjectID
	fileObjectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return &InvalidInputError{"Invalid file ID format"}
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return &InvalidInputError{"Invalid user ID format"}
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call repository with correct types
	return cs.Repo.CreateComment(ctx, fileObjectID, userObjectID, content)
}

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}
