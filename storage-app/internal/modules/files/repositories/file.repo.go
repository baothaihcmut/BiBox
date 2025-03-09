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
	ParentFolderId *primitive.ObjectID
	SortBy         string
	IsAsc          bool
	Offset         int
	Limit          int
	PermssionLimit int
	FileType       *enums.MimeType
	OwnerId        *primitive.ObjectID
	UserId         primitive.ObjectID
}

type FileRepository interface {
	CreateFile(context.Context, *models.File) error
	FindFileById(ctx context.Context, id primitive.ObjectID, isDeleted bool) (*models.File, error)
	UpdateFile(context.Context, *models.File) error
	FindFileWithPermssionAndCount(ctx context.Context, args FindFileOfUserArg) ([]*models.FileWithPermission, int64, error)
	GetSubFileRecursive(context.Context, primitive.ObjectID, []primitive.ObjectID) ([]*models.File, error)
}

type MongoFileRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

// GetSubFile implements FileRepository.
func (f *MongoFileRepository) GetSubFileRecursive(ctx context.Context, fileId primitive.ObjectID, excludeFileIds []primitive.ObjectID) ([]*models.File, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{{Key: "_id", Value: fileId}}},
		},
		bson.D{
			{Key: "$graphLookup", Value: bson.D{
				{Key: "from", Value: "files"},
				{Key: "startWith", Value: "$_id"}, // Use actual ObjectID
				{Key: "connectFromField", Value: "_id"},
				{Key: "connectToField", Value: "parent_folder_id"},
				{Key: "as", Value: "sub_files"},
				{Key: "restrictSearchWithMatch", Value: bson.D{
					{Key: "$nin", Value: excludeFileIds},
				}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "sub_files", Value: 1},
			}},
		},
	}

	cursor, err := f.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	// Define a result struct

	var subFiles []*struct {
		SubFiles []*models.File `bson:"sub_files"`
	}
	if err := cursor.All(ctx, &subFiles); err != nil {
		return nil, err
	}

	// Check if result is empty
	if len(subFiles) == 0 || len(subFiles[0].SubFiles) == 0 {
		return nil, nil
	}

	fmt.Println(subFiles)
	return subFiles[0].SubFiles, nil
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

func (f *MongoFileRepository) UpdateFile(ctx context.Context, file *models.File) error {

	_, err := f.collection.UpdateOne(ctx, bson.M{
		"_id": file.ID,
	}, bson.M{
		"$set": file,
	})
	if err != nil {
		f.logger.Errorf(ctx, map[string]any{
			"file_id": file.ID,
		}, "Error update file document:", err)
		return err
	}
	return nil
}

func (f *MongoFileRepository) FindFileWithPermssionAndCount(ctx context.Context, args FindFileOfUserArg) ([]*models.FileWithPermission, int64, error) {
	filter := bson.D{}
	if args.OwnerId != nil {
		filter = append(filter, bson.E{Key: "owner_id", Value: *args.OwnerId})
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

			//lookup permission
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
					}},
				}},
			},
			//match file user have permission
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "permissions", Value: bson.D{
						{Key: "$elemMatch", Value: bson.D{
							{Key: "user_id", Value: args.UserId},
						}},
					}},
				}},
			},

			//add permission type of user
			bson.D{
				{Key: "$addFields", Value: bson.D{
					{Key: "permission_users", Value: bson.D{
						{Key: "$filter", Value: bson.D{
							{Key: "input", Value: "$permissions"},
							{Key: "as", Value: "item"},
							{Key: "cond", Value: bson.D{
								{Key: "$eq", Value: bson.A{"$$item.user_id", args.OwnerId}},
							}},
						}},
					}},
				}},
			},
			bson.D{
				{Key: "$addFields", Value: bson.D{
					{Key: "permission_type", Value: bson.D{
						{Key: "$arrayElemAt", Value: bson.A{"$permission_users.permission_type", 0}},
					}},
				}},
			},
			bson.D{
				{Key: "$unset", Value: "permission_types"},
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
						{Key: "$map", Value: bson.D{
							{Key: "input", Value: bson.D{
								{Key: "$slice", Value: bson.A{"$permissions", 0, args.PermssionLimit}},
							}},
							{Key: "as", Value: "perm"},
							{Key: "in", Value: bson.D{
								{Key: "user_id", Value: "$$perm.user_id"},
								{Key: "permission_type", Value: "$$perm.permission_type"},
								{Key: "user_email", Value: "$$perm.user.email"},
								{Key: "user_first_name", Value: "$$perm.user.first_name"},
								{Key: "user_last_name", Value: "$$perm.user.last_name"},
								{Key: "user_image", Value: "$$perm.user.image"},
							}},
						}},
					}},
					{Key: "permission_type", Value: 1},
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
