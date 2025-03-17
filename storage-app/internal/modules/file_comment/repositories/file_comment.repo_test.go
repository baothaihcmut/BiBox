package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestFetchCommentsWithUsersAndAnswers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Client.Disconnect(context.Background())

	mt.Run("fetch comments with users and answers", func(mt *mtest.T) {
		// Setup collections
		commentCollection := mt.Coll
		userCollection := mt.DB.Collection("users")
		answerCollection := mt.DB.Collection("answers")

		// Initialize repository
		repo := &CommentRepository{
			Collection:       commentCollection,
			UserCollection:   userCollection,
			AnswerCollection: answerCollection,
		}

		// Insert sample data
		userID := primitive.NewObjectID()
		answerID := primitive.NewObjectID()
		commentID := primitive.NewObjectID()

		mt.AddMockResponses(
			mtest.CreateSuccessResponse(),
			mtest.CreateSuccessResponse(),
			mtest.CreateSuccessResponse(),
		)

		_, err := userCollection.InsertOne(context.TODO(), bson.M{
			"_id":   userID,
			"email": "test@example.com",
			"name":  "Test User",
		})
		assert.NoError(t, err)

		_, err = answerCollection.InsertOne(context.TODO(), bson.M{
			"_id":         answerID,
			"user_id":     userID,
			"content":     "This is an answer",
			"answered_at": time.Now(),
		})
		assert.NoError(t, err)

		_, err = commentCollection.InsertOne(context.TODO(), bson.M{
			"_id":        commentID,
			"file_id":    primitive.NewObjectID(),
			"user_id":    userID,
			"comment":    "This is a comment",
			"created_at": time.Now(),
			"answers":    []primitive.ObjectID{answerID},
		})
		assert.NoError(t, err)

		// Mock the aggregation pipeline response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.comments", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: commentID},
			{Key: "file_id", Value: primitive.NewObjectID()},
			{Key: "user_id", Value: userID},
			{Key: "comment", Value: "This is a comment"},
			{Key: "created_at", Value: time.Now()},
			{Key: "answers", Value: []bson.M{
				{
					"_id":         answerID,
					"user_id":     userID,
					"content":     "This is an answer",
					"answered_at": time.Now(),
				},
			}},
			{Key: "user", Value: bson.M{
				"_id":   userID,
				"email": "test@example.com",
				"name":  "Test User",
			}},
		}))

		// Fetch comments with users and answers
		comments, err := repo.FetchCommentsWithUsersAndAnswers(context.TODO())
		assert.NoError(t, err)
		assert.Len(t, comments, 1)

		comment := comments[0]
		assert.Equal(t, commentID, comment["_id"])
		assert.Equal(t, userID, comment["user_id"])
		assert.Equal(t, "This is a comment", comment["comment"])

		user := comment["user"].(bson.M)
		assert.Equal(t, userID, user["_id"])
		assert.Equal(t, "test@example.com", user["email"])
		assert.Equal(t, "Test User", user["name"])

		answers := comment["answers"].([]interface{})
		assert.Len(t, answers, 1)

		answer := answers[0].(bson.M)
		assert.Equal(t, answerID, answer["_id"])
		assert.Equal(t, userID, answer["user_id"])
		assert.Equal(t, "This is an answer", answer["content"])
	})
}
