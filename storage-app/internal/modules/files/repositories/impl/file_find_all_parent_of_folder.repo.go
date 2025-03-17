package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (f *MongoFileRepository) FindAllParentFolder(ctx context.Context, fileId primitive.ObjectID) ([]*models.File, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: fileId},
			}},
		},
		bson.D{
			{Key: "$graphLookup", Value: bson.D{
				{Key: "from", Value: "files"},
				{Key: "startWith", Value: "$parent_folder_id"}, // Use actual ObjectID
				{Key: "connectFromField", Value: "parent_folder_id"},
				{Key: "connectToField", Value: "_id"},
				{Key: "as", Value: "parent_folders"},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "parent_folders", Value: 1},
			}},
		},
	}
	cursor, err := f.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var res []*struct {
		ParentFolders []*models.File `bson:"parent_folders"`
	}
	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	if len(res) == 0 || len(res[0].ParentFolders) == 0 {
		return nil, nil
	}
	return res[0].ParentFolders, nil

}
