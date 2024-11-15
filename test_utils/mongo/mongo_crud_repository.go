package mongo

import (
	"context"
	"github.com/bloock/go-kit/client"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoCrudRepository struct {
	client     *client.MongoClient
	collection string
}

func NewMongoCrudRepository(client *client.MongoClient, collection string) MongoCrudRepository {
	return MongoCrudRepository{
		client:     client,
		collection: collection,
	}
}

func (m MongoCrudRepository) Truncate() error {
	ctxWithTimeout, cancel := m.client.ContextWithTimeout(context.Background())
	defer cancel()

	collection := m.client.DB().Collection(m.collection)

	if _, err := collection.DeleteMany(ctxWithTimeout, bson.M{}); err != nil {
		return err
	}
	return nil
}
