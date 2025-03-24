package repositories

import (
	"context"
	"log"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository struct {
	Collection *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) *CommentRepository {
	return &CommentRepository{
		Collection: db.Collection("file_comments"),
	}
}

// retrieves all comments
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

// retrieves all comments with details
func (cr *CommentRepository) FetchCommentsWithUsersAndAnswers(ctx context.Context) ([]map[string]interface{}, error) {
	pipeline := mongo.Pipeline{
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "localField", Value: "user_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "user"},
			}},
		},
		{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$user"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
	}

	cursor, err := cr.Collection.Aggregate(ctx, pipeline)
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

// inserts a new comment
func (cr *CommentRepository) CreateComment(ctx context.Context, fileID, userID primitive.ObjectID, commentText string) error {
	comment := models.FileComment{
		FileID:    fileID,
		UserID:    userID,
		Content:   commentText,
		CreatedAt: time.Now(),
		Answers:   []models.AnswerComment{},
	}

	_, err := cr.Collection.InsertOne(ctx, comment)
	if err != nil {
		log.Println("Error inserting comment:", err)
		return err
	}

	log.Printf("Inserted comment: fileID=%s, userID=%s, comment=%s", fileID.Hex(), userID.Hex(), commentText)
	return nil
}

// retrieves comments
func (cr *CommentRepository) GetCommentsByFile(fileID string) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func (cr *CommentRepository) AnswerComment(ctx context.Context, commentID, userID primitive.ObjectID, content string) error {
	answer := models.AnswerComment{
		CommentID: commentID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	filter := bson.M{"_id": commentID}
	update := bson.M{
		"$push": bson.M{
			"answers": answer,
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

// retrieves user
func (cr *CommentRepository) FetchUserByID(ctx context.Context, userID primitive.ObjectID) (map[string]interface{}, error) {
	var user bson.M
	err := cr.Collection.Database().Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// retrieves answer
func (cr *CommentRepository) FetchAnswerByID(ctx context.Context, answerID primitive.ObjectID) (map[string]interface{}, error) {
	var answer bson.M
	err := cr.Collection.FindOne(ctx, bson.M{"answers._id": answerID}).Decode(&answer)
	if err != nil {
		return nil, err
	}
	return answer, nil
}
