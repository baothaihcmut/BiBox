package repositories

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilterPermssionOption int

const (
	PermssionInList FilterPermssionOption = iota
	PermssionLessThan
	PermssionGreaterThan
	PermssionEqual
	PermssionNotEqual
	PermssionNotInList
)

type FilterPermssionType struct {
	Option FilterPermssionOption
	Value  []enums.FilePermissionType
}

type FilePermissionRepository interface {
	UpdatePermission(ctx context.Context, fileID primitive.ObjectID, userID primitive.ObjectID, permissionType enums.FilePermissionType, accessSecure bool, canShare bool) error
	GetFileByID(ctx context.Context, fileID primitive.ObjectID) (*models.FilePermission, error)
	CreateFilePermission(ctx context.Context, filePermission *models.FilePermission) error
	GetFilePermission(ctx context.Context, fileID primitive.ObjectID, userID primitive.ObjectID, option FilterPermssionType) (*models.FilePermission, error)
}

type FilePermissionRepositoryImpl struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewPermissionRepository(collection *mongo.Collection, logger logger.Logger) FilePermissionRepository {
	return &FilePermissionRepositoryImpl{
		collection: collection,
		logger:     logger,
	}
}

// update permissino
func (pr *FilePermissionRepositoryImpl) UpdatePermission(
	ctx context.Context,
	fileID primitive.ObjectID,
	userID primitive.ObjectID,
	permissionType enums.FilePermissionType,
	accessSecure bool,
	canShare bool,
) error {
	filter := bson.M{"file_id": fileID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"permission_type":    permissionType,
			"access_secure_file": accessSecure,
			"can_share":          canShare,
		},
	}

	_, err := pr.collection.UpdateOne(ctx, filter, update)
	return err
}

// get file by ID to check ownership
func (pr *FilePermissionRepositoryImpl) GetFileByID(ctx context.Context, fileID primitive.ObjectID) (*models.FilePermission, error) {
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
func (pr *FilePermissionRepositoryImpl) CreateFilePermission(ctx context.Context, filePermission *models.FilePermission) error {
	_, err := pr.collection.InsertOne(ctx, filePermission)
	if err != nil {
		pr.logger.Errorf(ctx, map[string]any{
			"file_id": filePermission.FileID,
			"user_id": filePermission.UserID,
		}, "Error when insert file permission: ", err)
		return err
	}
	return nil
}

// get file permission by fileID and userID and permssion type
func (pr *FilePermissionRepositoryImpl) GetFilePermission(ctx context.Context, fileID primitive.ObjectID, userID primitive.ObjectID, option FilterPermssionType) (*models.FilePermission, error) {
	// build filter
	filterPermssionType := bson.M{}
	switch option.Option {
	case PermssionInList:
		filterPermssionType = bson.M{"$in": option.Value}
	case PermssionLessThan:
		filterPermssionType = bson.M{"$lt": option.Value[0]}
	case PermssionGreaterThan:
		filterPermssionType = bson.M{"$gt": option.Value[0]}
	case PermssionEqual:
		filterPermssionType = bson.M{"$eq": option.Value[0]}
	case PermssionNotEqual:
		filterPermssionType = bson.M{"$ne": option.Value[0]}
	case PermssionNotInList:
		filterPermssionType = bson.M{"$nin": option.Value[0]}
	}

	var file models.FilePermission
	err := pr.collection.FindOne(ctx, bson.M{
		"file_id":         fileID,
		"user_id":         userID,
		"permission_type": filterPermssionType,
	}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // File permission not found
		}
		return nil, err
	}
	return &file, nil
}
