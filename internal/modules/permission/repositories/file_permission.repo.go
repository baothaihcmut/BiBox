package repositories

import (
	"context"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PermissionRepository handles database operations
type PermissionRepository struct {
	collection *mongo.Collection
}

// NewPermissionRepository initializes a new repository
func NewPermissionRepository(db *mongo.Database) *PermissionRepository {
	return &PermissionRepository{
		collection: db.Collection("permissions"),
	}
}

// FetchPermissions retrieves all permissions from the database
func (pr *PermissionRepository) FetchPermissions() ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query all documents
	cursor, err := pr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []map[string]interface{}
	for cursor.Next(ctx) {
		var permission bson.M
		if err := cursor.Decode(&permission); err != nil {
			log.Println("Error decoding permission:", err)
			continue
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

// CreatePermission inserts a new permission record into the database
func (pr *PermissionRepository) CreatePermission(fileID primitive.ObjectID, userID primitive.ObjectID, permissionType int, accessSecure bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert into MongoDB
	_, err := pr.collection.InsertOne(ctx, bson.M{
		"file_id":            fileID,
		"user_id":            userID,
		"permission_type":    permissionType,
		"access_secure_file": accessSecure,
		"created_at":         time.Now(),
	})
	if err != nil {
		log.Println("Error inserting permission:", err)
		return err
	}

	log.Printf("Inserted permission: fileID=%s, userID=%s, type=%s, secure=%v", fileID, userID, strconv.Itoa(permissionType), accessSecure)

	return nil
}

// GetPermissionsByUser retrieves permissions for a specific user
func (pr *PermissionRepository) GetPermissionsByUser(userID string) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := pr.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []map[string]interface{}
	for cursor.Next(ctx) {
		var permission bson.M
		if err := cursor.Decode(&permission); err != nil {
			log.Println("Error decoding permission:", err)
			continue
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}
