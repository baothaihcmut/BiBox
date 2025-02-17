package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoService interface {
	BeginTransaction(context.Context) (mongo.Session, error)
	EndTransansaction(context.Context, mongo.Session)
	CommitTransaction(context.Context, mongo.Session) error
	RollbackTransaction(context.Context, mongo.Session) error
}

type MongoTransactionServiceImpl struct {
	mongoClient *mongo.Client
}

// EndTransansaction implements MongoTransactionService.
func (m *MongoTransactionServiceImpl) EndTransansaction(ctx context.Context, session mongo.Session) {
	session.EndSession(ctx)
}

// BeginTransaction implements MongoTransactionService.
func (m *MongoTransactionServiceImpl) BeginTransaction(_ context.Context) (mongo.Session, error) {
	session, err := m.mongoClient.StartSession()
	if err != nil {
		return nil, err
	}
	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	return session, nil
}

// CommitTransaction implements MongoTransactionService.
func (m *MongoTransactionServiceImpl) CommitTransaction(ctx context.Context, session mongo.Session) error {
	return session.CommitTransaction(ctx)
}

// RollbackTransaction implements MongoTransactionService.
func (m *MongoTransactionServiceImpl) RollbackTransaction(ctx context.Context, session mongo.Session) error {
	return session.AbortTransaction(ctx)
}

func NewMongoTransactionService(mongoClient *mongo.Client) MongoService {
	return &MongoTransactionServiceImpl{
		mongoClient: mongoClient,
	}
}
