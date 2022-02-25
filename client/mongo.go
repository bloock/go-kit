package client

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient struct {
	client *mongoDriver.Client
	logger zerolog.Logger
}

func NewMongoClient(user, pass, host, port string, isCosmos bool, timeout time.Duration, l zerolog.Logger) (*MongoClient, error) {
	l = l.With().Str("layer", "infrastructure").Str("component", "mongo").Logger()

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?retrywrites=false&maxIdleTimeMS=120000", user, pass, host, port)
	if isCosmos {
		mongoURI = mongoURI + "&ssl=true&replicaSet=globaldb"
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongoDriver.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &MongoClient{
		client: client,
		logger: l,
	}, nil
}

func (c MongoClient) DB() *mongoDriver.Client {
	return c.client
}

func (c MongoClient) MigrateUp(dbName string, path string) error {
	driver, err := mongodb.WithInstance(c.DB(), &mongodb.Config{
		DatabaseName:    dbName,
		TransactionMode: false,
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		dbName,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}

		return err
	}

	return nil
}

func (c MongoClient) MigrateDown(dbName string, path string) error {
	driver, err := mongodb.WithInstance(c.DB(), &mongodb.Config{
		DatabaseName:    dbName,
		TransactionMode: false,
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		dbName,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil {
		return err
	}

	return nil
}
