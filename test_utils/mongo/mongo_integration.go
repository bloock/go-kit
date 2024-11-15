package mongo

import (
	"context"
	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/observability"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"testing"
)

const (
	containerName = "mongo_integration_test"
	imageName     = "mongo"
	imageTag      = "latest"
	dbName        = "test"
)

var mongoClient *client.MongoClient

func SetupMongoIntegrationTest(m *testing.M, testTimeout uint, migrationPath ...string) {
	pool, resource := initDB(testTimeout, migrationPath...)
	code := m.Run()
	closeDB(pool, resource)
	os.Exit(code)
}

func initDB(testTimeout uint, migrationPath ...string) (*dockertest.Pool, *dockertest.Resource) {
	ctx := context.Background()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	if err = pool.RemoveContainerByName(containerName); err != nil {
		log.Fatalf("%s", err)
	}

	opt := dockertest.RunOptions{
		Name:       containerName,
		Repository: imageName,
		Tag:        imageTag,
		Env: []string{
			"MONGO_INITDB_ROOT_USERNAME=test",
			"MONGO_INITDB_ROOT_PASSWORD=test",
		},
	}

	resource, err := pool.RunWithOptions(&opt)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	err = resource.Expire(testTimeout)
	if err != nil {
		log.Fatalf("Error setting expiration: %s", err)
	}

	if err = pool.Retry(func() error {
		mongoClient, err = client.NewMongoClient("test", "test", "localhost",
			resource.GetPort("27017/tcp"), dbName, false, observability.Logger{})
		if err != nil {
			return err
		}

		return mongoClient.Client().Ping(ctx, nil)
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	for _, mp := range migrationPath {
		if err = mongoClient.MigrateUp(dbName, mp); err != nil {
			log.Fatalf("Error executing migration: %s", err.Error())
		}
	}
	return pool, resource
}

func closeDB(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func MongoClient() *client.MongoClient {
	return mongoClient
}
