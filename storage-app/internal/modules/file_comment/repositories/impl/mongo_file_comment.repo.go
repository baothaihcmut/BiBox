package impl

import "go.mongodb.org/mongo-driver/mongo"

type MongoFileCommentRepo struct {
	collection *mongo.Collection
}

func NewMongoFileCommentRepo(collection *mongo.Collection) *MongoFileCommentRepo {
	return &MongoFileCommentRepo{
		collection: collection,
	}
}
