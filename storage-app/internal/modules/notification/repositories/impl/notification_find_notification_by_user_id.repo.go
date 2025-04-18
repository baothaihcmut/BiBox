package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (n *NotificationRepoImpl) FindNotificationByUserIdSortByCreatedAtAndCount(
	ctx context.Context,
	useId primitive.ObjectID,
	limit int,
	offset int) ([]*models.Notification, int, error) {
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 1)
	doneCh := make(chan struct{}, 1)
	dataCh := make(chan []*models.Notification, 1)
	countCh := make(chan int, 1)
	filter := bson.M{
		"user_id": useId,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		findOptions := options.Find().
			SetSort(bson.D{{Key: "created_at", Value: -1}}).
			SetSkip(int64(offset)).
			SetLimit(int64(limit))
		cursor, err := n.collection.Find(ctx, filter, findOptions)
		if err != nil {
			errCh <- err
			return
		}
		defer cursor.Close(ctx)
		var res []*models.Notification
		if err := cursor.All(ctx, &res); err != nil {
			errCh <- err
			return
		}
		dataCh <- res
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := n.collection.CountDocuments(ctx, filter)
		if err != nil {
			errCh <- err
			return
		}
		countCh <- int(count)
	}()

	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()
	select {
	case err := <-errCh:
		return nil, 0, err
	case <-doneCh:
		return <-dataCh, <-countCh, nil
	}
}
