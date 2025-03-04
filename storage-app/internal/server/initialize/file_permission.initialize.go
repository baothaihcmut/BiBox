package initialize

import (
	"context"
	"log"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitializePermissionModule initializes the permission module
func InitializePermissionModule(client *mongo.Client) (*repositories.PermissionRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database("storage-app")
	permissionRepo := repositories.NewPermissionRepository(db)

	// Ensure indexes
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"user_id": 1, "file_id": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := permissionRepo.Collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Println("Error creating index for permissions:", err)
		return nil, err
	}

	return permissionRepo, nil
}
