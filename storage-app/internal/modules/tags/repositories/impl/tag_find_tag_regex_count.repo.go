package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (t *MongoTagRepository) FindAllTagsAndCount(ctx context.Context, query string, limit, offset int) ([]*models.Tag, int, error) {

	filter := bson.M{
		"name": bson.M{
			"$regex":   query,
			"$options": "i",
		},
	}
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)
	doneCh := make(chan struct{}, 1)
	var data []*models.Tag
	var count int
	wg.Add(1)
	go func() {
		defer wg.Done()

		cursor, err := t.collection.Find(ctx, filter, options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)))
		if err != nil {
			errCh <- err
			return
		}
		if err := cursor.All(ctx, &data); err != nil {
			errCh <- err
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		countRes, err := t.collection.CountDocuments(ctx, filter)
		if err != nil {
			errCh <- err
			return
		}
		count = int(countRes)
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
