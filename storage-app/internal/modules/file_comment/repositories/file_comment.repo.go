package repositories

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CommentRepo
type CommentRepository struct {
	collection *mongo.Collection
}

// NewCommentRepos
func NewCommentRepository(db *mongo.Database) *CommentRepository {
	return &CommentRepository{
		collection: db.Collection("comments"),
	}
}

// FetchComments retrieves all comments from the database
func (cr *CommentRepository) FetchComments() ([]map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := cr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []map[string]any
	for cursor.Next(ctx) {
		var comment bson.M
		if err := cursor.Decode(&comment); err != nil {
			log.Println("Error decoding comment:", err)
			continue
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// Fix: CreateComment now accepts `primitive.ObjectID` for fileID and userID
func (cr *CommentRepository) CreateComment(ctx context.Context, fileID, userID primitive.ObjectID, commentText string) error {
	_, err := cr.collection.InsertOne(ctx, bson.M{
		"file_id":    fileID,
		"user_id":    userID,
		"comment":    commentText,
		"created_at": time.Now(),
	})
	if err != nil {
		log.Println("Error inserting comment:", err)
		return err
	}

	log.Printf("Inserted comment: fileID=%s, userID=%s, comment=%s", fileID.Hex(), userID.Hex(), commentText)
	return nil
}

// GetCommentsByFile retrieves comments for a specific file
func (cr *CommentRepository) GetCommentsByFile(fileID string) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert fileID from string to ObjectID
	fileObjectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, err // Invalid fileID format
	}

	cursor, err := cr.collection.Find(ctx, bson.M{"file_id": fileObjectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []map[string]any
	for cursor.Next(ctx) {
		var comment bson.M
		if err := cursor.Decode(&comment); err != nil {
			log.Println("Error decoding comment:", err)
			continue
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
