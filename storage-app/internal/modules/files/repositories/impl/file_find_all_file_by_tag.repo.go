package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (f *MongoFileRepository) FindAllFileByTagAndCount(
	ctx context.Context,
	tagId, userId primitive.ObjectID,
	limit, offset int,
	sortBy string,
	isAsc bool,
) ([]*models.FileWithPermission, int, error) {
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)
	doneCh := make(chan struct{}, 1)
	var data []*models.FileWithPermission
	var count int
	wg.Add(1)
	go func() {
		defer wg.Done()
		pipeline := mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: bson.M{
					"tag_ids": bson.M{
						"$elemMatch": bson.M{
							"$eq": tagId,
						},
					},
				}},
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
							{Key: "user_id", Value: userId},
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
								{Key: "$eq", Value: bson.A{"$$item.user_id", userId}},
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
								{Key: "$slice", Value: bson.A{"$permissions", 0, 4}},
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
					{Key: sortBy, Value: 1},
				}},
			},
			bson.D{
				{Key: "$skip", Value: offset},
			},
			bson.D{
				{Key: "$limit", Value: limit},
			},
		}
		cursor, err := f.collection.Aggregate(ctx, pipeline)
		if err != nil {
			errCh <- err
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.All(ctx, &data); err != nil {
			errCh <- err
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		pipeline := mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: bson.M{
					"tag_ids": bson.M{
						"$elemMatch": bson.M{
							"$eq": tagId,
						},
					},
				}},
			},
			bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "file_permissions"},
					{Key: "localField", Value: "_id"},
					{Key: "foreignField", Value: "file_id"},
					{Key: "as", Value: "permissions"},
				}},
			},
			//match file user have permission
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "permissions", Value: bson.D{
						{Key: "$elemMatch", Value: bson.D{
							{Key: "user_id", Value: userId},
						}},
					}},
				}},
			},
			bson.D{
				{Key: "$count", Value: "count_total"},
			},
		}
		cursor, err := f.collection.Aggregate(ctx, pipeline)
		if err != nil {
			errCh <- err
			return
		}
		defer cursor.Close(ctx)
		var res struct {
			TotalCount int `bson:"count_total"`
		}
		if cursor.Next(ctx) {
			if err := cursor.Decode(&res); err != nil {
				errCh <- err
			}
			count = res.TotalCount
		}
	}()
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()
	select {
	case <-doneCh:
		return data, count, nil
	case err := <-errCh:
		return nil, 0, err
	}
}
