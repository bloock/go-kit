package client

import (
	"context"
	"fmt"
	"time"

	"github.com/bloock/go-kit/observability"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	defaultTimeout = 5 * time.Second
)

type MongoClient struct {
	client       *mongoDriver.Client
	databaseName string
	timeout      time.Duration
	logger       observability.Logger
}

func NewMongoClient(user, pass, host, port, databaseName string, isCosmos bool, l observability.Logger, opts ...ClientOpt) (*MongoClient, error) {
	l.UpdateLogger(l.With().Str("layer", "infrastructure").Str("component", "mongo").Logger())

	op := &clientOpts{
		timeout:  defaultTimeout,
		readPref: readpref.Nearest(),
	}

	for _, fn := range opts {
		fn(op)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s/?retrywrites=false&maxIdleTimeMS=120000", host, port)
	if isCosmos {
		mongoURI = mongoURI + "&ssl=true&replicaSet=globaldb"
	}

	ctx, cancel := context.WithTimeout(context.Background(), op.timeout)
	defer cancel()

	client, err := mongoDriver.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &MongoClient{
		client:       client,
		databaseName: databaseName,
		timeout:      op.timeout,
		logger:       l,
	}, nil
}

func (c MongoClient) Client() *mongoDriver.Client {
	return c.client
}

func (c MongoClient) DB() *mongoDriver.Database {
	return c.client.Database(c.databaseName)
}

func (c MongoClient) MigrateUp(dbName string, path string) error {
	driver, err := mongodb.WithInstance(c.Client(), &mongodb.Config{
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

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}

		return err
	}

	return nil
}

func (c MongoClient) MigrateDown(dbName string, path string) error {
	driver, err := mongodb.WithInstance(c.Client(), &mongodb.Config{
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

type ClientOpt func(opts *clientOpts)

type clientOpts struct {
	timeout  time.Duration
	readPref *readpref.ReadPref
}

func WithTimeout(timeout time.Duration) ClientOpt {
	return func(opts *clientOpts) {
		opts.timeout = timeout
	}
}

func WithReadPref(pref *readpref.ReadPref) ClientOpt {
	return func(opts *clientOpts) {
		opts.readPref = pref
	}
}
