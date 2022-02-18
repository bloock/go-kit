package health

import (
	"context"

	"github.com/bloock/go-kit/client"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type HealthMongo struct {
	client      *client.MongoClient
	description string
	version     string
}

func NewHealthMongo(client *client.MongoClient, description string) HealthMongo {
	return HealthMongo{
		client:      client,
		description: description,
		version:     DepVersion("mongo-driver"),
	}
}

func (h HealthMongo) HealthCheck() ExternalServiceDetails {
	s := "pass"
	var e string

	ctx := context.Background()
	err := h.client.DB().Ping(ctx, readpref.Primary())
	if err != nil {
		s = "error"
		e = err.Error()
	}

	return ExternalServiceDetails{
		Description: h.description,
		Version:     h.version,
		Status:      s,
		Error:       e,
	}
}
