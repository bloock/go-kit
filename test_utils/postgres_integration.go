package test_utils

import (
	"context"
	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/observability"
	"github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"runtime"
	"testing"
	"time"
)

const (
	containerName = "postgres_integration_test"
	imageName     = "postgres"
	imageTag      = "11"
)

type Logger struct {
}

func (l Logger) Print(v ...interface{}) {
	for _, e := range v {
		err := e.(error)
		if err.Error() == "unexpected EOF" {
			return
		}
	}
	log.Println(v...)
}

var postgresSQLClient *client.PostgresSQLClient

func SetupPostgresIntegrationTest(m *testing.M, testTimeout uint, migrationPath ...string) {
	pool, resource := initDB(testTimeout, migrationPath...)
	code := m.Run()
	closeDB(pool, resource)
	os.Exit(code)
}

func initDB(testTimeout uint, migrationPath ...string) (*dockertest.Pool, *dockertest.Resource) {
	ctx := context.Background()

	var platform string
	mysql.SetLogger(Logger{})

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	if err = pool.RemoveContainerByName(containerName); err != nil {
		log.Fatalf("%s", err)
	}

	if runtime.GOARCH == "arm64" {
		platform = "linux/x86-64"
	}

	opt := dockertest.RunOptions{
		Name:       containerName,
		Repository: imageName,
		Tag:        imageTag,
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
		},
		Platform: platform,
	}

	resource, err := pool.RunWithOptions(&opt)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	resource.Expire(testTimeout)

	if err := pool.Retry(func() error {
		postgresSQLClient, err = client.NewPostgresClient(ctx, "test", "test", "localhost",
			resource.GetPort("3306/tcp"), "test", false, &client.SQLConnOpts{
				MaxConnLifeTime: 5 * time.Minute,
				MaxOpenConns:    1,
				MaxIdleConns:    1,
			}, observability.Logger{})
		if err != nil {
			return err
		}

		return postgresSQLClient.DB().Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	for _, mp := range migrationPath {
		if err = postgresSQLClient.MigrateUp(ctx, mp); err != nil {
			log.Fatalf("%s", err.Error())
		}
	}

	return pool, resource
}

func closeDB(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func GetMysqlClient() *client.PostgresSQLClient {
	return postgresSQLClient
}
