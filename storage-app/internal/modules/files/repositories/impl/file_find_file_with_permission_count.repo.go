package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (f *MongoFileRepository) FindFileWithPermssionAndCount(ctx context.Context, args repositories.FindFileWithPermissionArg) ([]*models.FileWithPermission, int64, error) {
	filter := bson.D{}
	if args.OwnerId != nil {
		filter = append(filter, bson.E{Key: "owner_id", Value: args.OwnerId})
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
	if args.IsDeleted != nil {
		filter = append(filter, bson.E{Key: "is_deleted", Value: *args.IsDeleted})
	}
	if args.TagId != nil {
		filter = append(filter, bson.E{Key: "tag_ids", Value: bson.M{
			"$elemMatch": bson.M{
				"$eq": args.TagId,
			},
		}})
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
					{Key: "total_size", Value: 1},
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
								{Key: "$slice", Value: bson.A{"$permissions", 0, 3}},
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
