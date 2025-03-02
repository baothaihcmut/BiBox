package repositories

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PermissionRepository struct {
	collection *mongo.Collection
}

func NewPermissionRepository(db *mongo.Database) *PermissionRepository {
	return &PermissionRepository{
		collection: db.Collection("permissions"),
	}
}

// update permissino
func (pr *PermissionRepository) UpdatePermission(ctx context.Context, fileID primitive.ObjectID, userID primitive.ObjectID, permissionType int, accessSecure bool) error {
	filter := bson.M{"file_id": fileID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"permission_type":    permissionType,
			"access_secure_file": accessSecure,
		},
	}

	_, err := pr.collection.UpdateOne(ctx, filter, update)
	return err
}

// get file by ID to check ownership
func (pr *PermissionRepository) GetFileByID(ctx context.Context, fileID primitive.ObjectID) (*models.FilePermission, error) {
	var file models.FilePermission
	err := pr.collection.FindOne(ctx, bson.M{"file_id": fileID}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // File not found
		}
		return nil, err
	}
	return &file, nil
}

// insert file permission into DB
func (pr *PermissionRepository) CreateFilePermission(ctx context.Context, fileID primitive.ObjectID, ownerUserID string, canShare bool) error {
	permission := bson.M{
		"file_id":       fileID,
		"owner_user_id": ownerUserID,
		"can_share":     canShare,
	}

	_, err := pr.collection.InsertOne(ctx, permission)
	return err
}
func (pr *PermissionRepository) CheckUserPermission(ctx context.Context, fileID, userID primitive.ObjectID, allowedPermissions []int) (bool, error) {
	filter := bson.M{
		"file_id":         fileID,
		"user_id":         userID,
		"permission_type": bson.M{"$in": allowedPermissions},
	}

	count, err := pr.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
