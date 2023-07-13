package mysql

import (
	"context"
	"log"
	"os"
	"runtime"
	"testing"

	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/observability"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3"
)

const (
	containerName = "mysql_integration_test"
	imageName     = "mysql"
	imageTag      = "8.0.22"
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
	log.Println(v)
}

var mysqlClient *client.MysqlClient

func SetupMysqlIntegrationTest(m *testing.M, testTimeout uint, migrationPath ...string) {
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
			"MYSQL_ROOT_PASSWORD=test",
			"MYSQL_DATABASE=test",
		},
		Platform: platform,
	}

	resource, err := pool.RunWithOptions(&opt)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	resource.Expire(testTimeout)

	if err := pool.Retry(func() error {
		mysqlClient, err = client.NewMysqlClient(ctx, "root", "test", "localhost",
			resource.GetPort("3306/tcp"), "test", false, &client.SQLConnOpts{}, observability.Logger{})
		if err != nil {
			return err
		}

		return mysqlClient.DB().Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	for _, mp := range migrationPath {
		if err = mysqlClient.MigrateUp(ctx, mp); err != nil {
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

func GetMysqlClient() *client.MysqlClient {
	return mysqlClient
}
