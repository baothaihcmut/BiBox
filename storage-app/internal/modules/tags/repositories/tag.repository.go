package repositories

import (
	"context"
	"errors"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TagRepository interface {
	CreateTag(context.Context, *models.Tag) error
	BulkCreateTag(context.Context, []*models.Tag) error
	FindTagById(context.Context, primitive.ObjectID) (*models.Tag, error)
}

type MongoTagRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func (m *MongoTagRepository) FindTagById(ctx context.Context, id primitive.ObjectID) (*models.Tag, error) {
	var res models.Tag
	err := m.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		m.logger.Errorf(ctx, map[string]interface{}{
			"tag_id": id.Hex(),
		}, "Error find tag by id:", err)
		return nil, err
	}
	return &res, nil
}

func (m *MongoTagRepository) CreateTag(ctx context.Context, tag *models.Tag) error {
	_, err := m.collection.InsertOne(ctx, tag)
	if err != nil {
		m.logger.Errorf(ctx, map[string]interface{}{
			"tag_id": tag.ID,
		}, "Error insert into tag collection: ", err)
		return err
	}
	return nil
}
func (m *MongoTagRepository) BulkCreateTag(ctx context.Context, tags []*models.Tag) error {
	input := make([]interface{}, len(tags))
	for idx, tag := range tags {
		input[idx] = tag
	}
	_, err := m.collection.InsertMany(ctx, input)
	if err != nil {
		m.logger.Errorf(ctx, nil, "Error insert many into tag collection: ", err)
		return err
	}
	return nil
}

func NewMongoTagRepository(collection *mongo.Collection, logger logger.Logger) TagRepository {
	return &MongoTagRepository{
		collection: collection,
		logger:     logger,
	}
}
