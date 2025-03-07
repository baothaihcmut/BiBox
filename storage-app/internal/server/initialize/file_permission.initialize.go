package initialize

import (
	"context"
	"log"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitializePermissionModule initializes the permission module
func InitializePermissionModule(client *mongo.Client, logger logger.Logger) (repositories.FilePermissionRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database("storage-app")
	permissionRepo := repositories.NewPermissionRepository(db.Collection("permissions"), logger) // use NewPermissionRepository

	// Ensure indexes
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"user_id": 1, "file_id": 1},
		Options: options.Index().SetUnique(true),
	}
	collection := db.Collection("permissions")                // access collection directly
	_, err := collection.Indexes().CreateOne(ctx, indexModel) // use collection directly
	if err != nil {
		log.Println("Error creating index for permissions:", err)
		return nil, err
	}

	return permissionRepo, nil
}
