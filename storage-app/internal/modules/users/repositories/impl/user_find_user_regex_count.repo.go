package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoUserRepository) FindUserRegexAndCount(ctx context.Context, query string, limit, offset *int) ([]*models.User, int, error) {
	filter := bson.M{
		"$or": bson.A{
			bson.M{"email": bson.M{
				"$regex":   query,
				"$options": "i",
			}},
			bson.M{"first_name": bson.M{
				"$regex":   query,
				"$options": "i",
			}},
			bson.M{"last_name": bson.M{
				"$regex":   query,
				"$options": "i",
			}},
		},
	}
	var res []*models.User
	var count int
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		findOption := options.FindOptions{}
		if limit != nil {
			findOption.SetLimit(int64(*limit))
		}
		if offset != nil {
			findOption.SetSkip(int64(*offset))
		}
		cursor, err := m.collection.Find(ctx, filter, &findOption)
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
		if err := cursor.All(ctx, &res); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancel()
			errCh <- err

		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := m.collection.CountDocuments(ctx, filter)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancel()
			errCh <- err
		}
		count = int(res)

	}()
	wg.Wait()
	select {
	case err := <-errCh:
		return nil, 0, err
	default:
		return res, count, nil
	}
}
