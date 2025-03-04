package initialize

import (
	"context"
	"log"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitializeCommentModule initializes the comment module
func InitializeCommentModule(client *mongo.Client) (*repositories.CommentRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database("storage-app")
	commentRepo := repositories.NewCommentRepository(db)

	// Ensure indexes
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"file_id": 1},
		Options: options.Index().SetUnique(false),
	}
	_, err := commentRepo.Collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Println("Error creating index for comments:", err)
		return nil, err
	}

	return commentRepo, nil
}
