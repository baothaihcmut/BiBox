package repositories

import (
	"log"
)

type CommentRepository struct{}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

func (cr *CommentRepository) FetchComments() ([]map[string]interface{}, error) {
	comments := []map[string]interface{}{
		{"file_id": "123", "user_id": "abc", "content": "Great file!"},
	}
	return comments, nil
}

func (cr *CommentRepository) InsertComment(fileID, userID, content string) error {
	log.Printf("Inserted comment: fileID=%s, userID=%s, content=%s", fileID, userID, content)
	return nil
}
