package health

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type HealthMongo struct {
	client      *mongo.Client
	description string
	version     string
}

func NewHealthMongo(client *mongo.Client, description, version string) HealthMongo {
	return HealthMongo{
		client:      client,
		description: description,
		version:     version,
	}
}

func (h HealthMongo) HealthCheck() ExternalServiceDetails {
	s := "pass"
	var e string

	ctx := context.Background()
	err := h.client.Ping(ctx, readpref.Primary())
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
