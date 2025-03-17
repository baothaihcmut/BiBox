package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetSubFile implements FileRepository.
func (f *MongoFileRepository) FindSubFileRecursive(ctx context.Context, fileId primitive.ObjectID) ([]*models.File, error) {
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

	return subFiles[0].SubFiles, nil
}
