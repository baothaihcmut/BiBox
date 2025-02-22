package repositories

import (
	"context"

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
func (pr *PermissionRepository) UpdatePermission(ctx context.Context, fileID, userID primitive.ObjectID, permissionType int, accessSecure bool) error {
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
