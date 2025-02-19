package repositories

import (
	"context"
	"fmt"

	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	"github.com/baothaihcmut/Storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FindFileOfUserArg struct {
	IsFolder   *bool
	IsInFolder *bool
	SortBy     string
	IsAsc      bool
	Offset     int
	Limit      int
}

type FileRepository interface {
	CreateFile(context.Context, *models.File) error
	FindFileById(ctx context.Context, id primitive.ObjectID, isDeleted bool) (*models.File, error)
	UploadedFile(context.Context, *models.File) error
	FindAllFileOfUser(ctx context.Context, userId primitive.ObjectID, args FindFileOfUserArg) ([]*models.File, error)
}

type MongoFileRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func (f *MongoFileRepository) CreateFile(ctx context.Context, file *models.File) error {
	_, err := f.collection.InsertOne(ctx, file)
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error insert document:", err)
	}
	return nil
}

func (f *MongoFileRepository) FindFileById(ctx context.Context, id primitive.ObjectID, isDeleted bool) (*models.File, error) {
	var res models.File
	err := f.collection.FindOne(ctx, bson.M{
		"_id":        id,
		"is_deleted": isDeleted,
	}).Decode(&res)
	if err != nil {
		f.logger.Errorf(ctx, map[string]interface{}{
			"file_id":    id,
			"is_deleted": isDeleted,
		}, "Error find file by id:", err)
		return nil, err
	}
	return &res, nil
}

func (f *MongoFileRepository) UploadedFile(ctx context.Context, file *models.File) error {
	_, err := f.collection.UpdateOne(ctx, bson.M{
		"_id": file.ID,
	}, bson.M{
		"$set": bson.M{
			"storage_detail.is_uploaded": true,
		},
	})
	if err != nil {
		f.logger.Errorf(ctx, map[string]interface{}{
			"file_id": file.ID,
		}, "Error update file document:", err)
		return err
	}
	return nil
}

func (f *MongoFileRepository) FindAllFileOfUser(ctx context.Context, userId primitive.ObjectID, args FindFileOfUserArg) ([]*models.File, error) {
	//for filter
	filter := bson.M{
		"owner_id": userId,
	}
	if args.IsFolder != nil {
		filter["is_folder"] = *args.IsFolder
	}
	if args.IsInFolder != nil {
		if *args.IsInFolder {
			filter["parent_folder_id"] = bson.M{
				"$ne": nil,
			}
		} else {
			filter["parent_folder_id"] = nil
		}
	}
	fmt.Println(filter)
	sort := bson.D{{Key: args.SortBy}}
	if args.IsAsc {
		sort[0].Value = 1
	} else {
		sort[0].Value = -1
	}
	opts := options.Find().SetSort(sort).SetSkip(int64(args.Offset)).SetLimit(int64(args.Limit))
	cursor, err := f.collection.Find(ctx, filter, opts)
	if err != nil {
		f.logger.Errorf(ctx, map[string]interface{}{
			"user_id": userId.Hex(),
		}, "Error find all file of user: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var files []*models.File
	if err := cursor.All(ctx, &files); err != nil {
		f.logger.Errorf(ctx, nil, "Error decode array of file document: ", err)
		return nil, err
	}
	return files, nil
}

func NewMongoFileRepo(collection *mongo.Collection, logger logger.Logger) FileRepository {
	return &MongoFileRepository{
		collection: collection,
		logger:     logger,
	}
}
