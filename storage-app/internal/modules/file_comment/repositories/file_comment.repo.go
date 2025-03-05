package repositories

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CommentRepository
type CommentRepository struct {
	Collection *mongo.Collection
}

// NewCommentRepository
func NewCommentRepository(db *mongo.Database) *CommentRepository {
	return &CommentRepository{
		Collection: db.Collection("file_comments"),
	}
}

// FetchComments retrieves all comments from the database
func (cr *CommentRepository) FetchComments() ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := cr.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []map[string]interface{}
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

// CreateComment inserts a new comment into the database
func (cr *CommentRepository) CreateComment(ctx context.Context, fileID, userID primitive.ObjectID, commentText string) error {
	_, err := cr.Collection.InsertOne(ctx, bson.M{
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

// GetCommentsByFile retrieves comments by file ID
func (cr *CommentRepository) GetCommentsByFile(fileID string) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert fileID from string to ObjectID
	fileObjectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, err
	}

	cursor, err := cr.Collection.Find(ctx, bson.M{"file_id": fileObjectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []map[string]interface{}
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

// AnswerComment
func (cr *CommentRepository) AnswerComment(ctx context.Context, commentID, userID primitive.ObjectID, content string) error {
	filter := bson.M{"_id": commentID}
	update := bson.M{
		"$set": bson.M{
			"answer":      content,
			"answered_by": userID,
			"answered_at": time.Now(),
		},
	}

	result, err := cr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating comment with answer:", err)
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	log.Printf("Answered comment: commentID=%s, userID=%s, answer=%s", commentID.Hex(), userID.Hex(), content)
	return nil
}

//moi comment co 1 mang ID cua answer, mang answerID string thanh primitive.objectID
// mang phan tu co noi dung cau tra loi, userID, answerd at, bo email
