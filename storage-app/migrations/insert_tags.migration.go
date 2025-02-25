package migrations

import (
	"context"
	"fmt"
	"log"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tags = []interface{}{
	models.Tag{ID: primitive.NewObjectID(), Name: "Technology"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Health"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Finance"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Education"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Science"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Sports"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Travel"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Food"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Music"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Movies"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Art"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Business"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Gaming"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Politics"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Environment"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Culture"},
	models.Tag{ID: primitive.NewObjectID(), Name: "History"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Literature"},
	models.Tag{ID: primitive.NewObjectID(), Name: "Photography"},
}

func InsertTags(ctx context.Context, client *mongo.Client, database string) error {
	collection := client.Database(database).Collection("tags")
	for _, tag := range tags {
		// Check if the tag already exists based on the 'Name' field
		filter := bson.M{"name": tag.(models.Tag).Name}
		update := bson.M{
			"$setOnInsert": tag, // Only insert if it doesn't already exist
		}

		// Insert the tag if it doesn't already exist
		opts := options.Update().SetUpsert(true)
		_, err := collection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("Error inserting tag: %v", err)
			return err
		}
	}
	fmt.Println("Tags inserted successfully.")
	return nil
}
