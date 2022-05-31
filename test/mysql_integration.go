package test

import (
	"github.com/bloock/go-kit/client"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog"
	"log"
	"os"
	"runtime"
	"testing"
)

const (
	containerName = "mysql_integration_test"
	imageName     = "mysql"
	imageTag      = "8.0.22"
)

type Logger struct {
}

func (l Logger) Print(v ...interface{}) {
	log.Println(v)
}

var mysqlClient *client.MysqlClient

func SetupMysqlIntegrationTest(m *testing.M, migrationPath string, testTimeout uint) {
	pool, resource := initDB(migrationPath, testTimeout)
	code := m.Run()
	closeDB(pool, resource)
	os.Exit(code)
}

func initDB(migrationPath string, testTimeout uint) (*dockertest.Pool, *dockertest.Resource) {
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
			"MYSQL_USER=test",
			"MYSQL_PASSWORD=test",
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
		mysqlClient, err = client.NewMysqlClient("test", "test", "localhost",
			resource.GetPort("3306/tcp"), "test", zerolog.Logger{})
		if err != nil {
			return err
		}

		return mysqlClient.DB().Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	if err = mysqlClient.MigrateUp(migrationPath); err != nil {
		log.Fatalf("%s", err.Error())
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
