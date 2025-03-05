package repositories

import (
	"context"
	"fmt"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FindFileOfUserArg struct {
	IsFolder       *bool
	ParentFolderId *string
	SortBy         string
	IsAsc          bool
	Offset         int
	Limit          int
	PermssionLimit int
	FileType       *enums.MimeType
	OwnerId        primitive.ObjectID
}

type FileRepository interface {
	CreateFile(context.Context, *models.File) error
	FindFileById(ctx context.Context, id primitive.ObjectID, isDeleted bool) (*models.File, error)
	UploadedFile(context.Context, *models.File) error
	FindAllFileOfUserWithPermssionAndCount(ctx context.Context, userId primitive.ObjectID, args FindFileOfUserArg) ([]*models.FileWithPermission, int64, error)
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
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		f.logger.Errorf(ctx, map[string]any{
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
		f.logger.Errorf(ctx, map[string]any{
			"file_id": file.ID,
		}, "Error update file document:", err)
		return err
	}
	return nil
}

func (f *MongoFileRepository) FindAllFileOfUserWithPermssionAndCount(ctx context.Context, userId primitive.ObjectID, args FindFileOfUserArg) ([]*models.FileWithPermission, int64, error) {
	filter := bson.D{
		{Key: "owner_id", Value: args.OwnerId},
	}
	if args.IsFolder != nil {
		filter = append(filter, bson.E{Key: "is_folder", Value: *args.IsFolder})
	}
	if args.ParentFolderId != nil {
		filter = append(filter, bson.E{Key: "parent_folder_id", Value: *args.ParentFolderId})
	} else {
		filter = append(filter, bson.E{Key: "parent_folder_id", Value: nil})
	}
	if args.FileType != nil {
		filter = append(filter, bson.E{Key: "storage_detail.mime_type", Value: *args.FileType})
	}

	fmt.Println(filter)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dataCh := make(chan []*models.FileWithPermission, 1)
	countCh := make(chan int64, 1)
	errCh := make(chan error, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		pipeline := mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: filter},
			},
			bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "file_permissions"},
					{Key: "localField", Value: "_id"},
					{Key: "foreignField", Value: "file_id"},
					{Key: "as", Value: "permissions"},
					{Key: "pipeline", Value: mongo.Pipeline{
						bson.D{
							{Key: "$lookup", Value: bson.D{
								{Key: "from", Value: "users"},
								{Key: "localField", Value: "user_id"},
								{Key: "foreignField", Value: "_id"},
								{Key: "as", Value: "user"},
							}},
						},
						bson.D{
							{Key: "$unwind", Value: "$user"},
						},
						bson.D{
							{Key: "$project", Value: bson.D{
								{Key: "user_id", Value: 1},
								{Key: "permission_type", Value: 1},
								{Key: "user_image", Value: "$user.image"},
								{Key: "user_email", Value: "$user.email"},
								{Key: "user_first_name", Value: "$user.first_name"},
								{Key: "user_last_name", Value: "$user.last_name"},
							}},
						},
					}},
				}},
			},
			bson.D{
				{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 1},
					{Key: "name", Value: 1},
					{Key: "owner_id", Value: 1},
					{Key: "is_folder", Value: 1},
					{Key: "parent_folder_id", Value: 1},
					{Key: "created_at", Value: 1},
					{Key: "updated_at", Value: 1},
					{Key: "opened_at", Value: 1},
					{Key: "has_password", Value: 1},
					{Key: "password", Value: 1},
					{Key: "description", Value: 1},
					{Key: "is_deleted", Value: 1},
					{Key: "deleted_at", Value: 1},
					{Key: "is_secure", Value: 1},
					{Key: "tags", Value: 1},
					{Key: "storage_detail", Value: 1},
					{Key: "permissions", Value: bson.D{
						{Key: "$slice", Value: bson.A{"$permissions", 0, args.PermssionLimit}},
					}},
				}},
			},
			bson.D{
				{Key: "$sort", Value: bson.D{
					{Key: args.SortBy, Value: 1},
				}},
			},
			bson.D{
				{Key: "$skip", Value: args.Offset},
			},
			bson.D{
				{Key: "$limit", Value: args.Limit},
			},
		}
		cursor, err := f.collection.Aggregate(ctx, pipeline)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancel()
			errCh <- err
			return
		}
		defer cursor.Close(ctx)

		result := make([]*models.FileWithPermission, 0)
		if err := cursor.All(ctx, &result); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancel()
			errCh <- err
			return
		}
		dataCh <- result
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := f.collection.CountDocuments(ctx, filter)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancel()
			errCh <- err
			return
		}
		countCh <- count
	}()
	wg.Wait()
	select {
	case err := <-errCh:
		return nil, 0, err
	default:
		return <-dataCh, <-countCh, nil
	}
}

func NewMongoFileRepo(collection *mongo.Collection, logger logger.Logger) FileRepository {
	return &MongoFileRepository{
		collection: collection,
		logger:     logger,
	}
}
