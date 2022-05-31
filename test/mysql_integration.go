package test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3"
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

var db *sql.DB

func SetupMysqlIntegrationTest(m *testing.M, migrationPath string, testTimeout uint) {
	pool, resource := initDB(migrationPath, testTimeout)
	code := m.Run()
	closeDB(pool, resource)
	os.Exit(code)
}

func initDB(migrationPath string, testTimeout uint) (*dockertest.Pool, *dockertest.Resource) {
	var platform string

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
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
			"LOGGING_LEVEL=ERROR",
		},
		Platform: platform,
	}

	resource, err := pool.RunWithOptions(&opt)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	resource.Expire(testTimeout)

	if err := pool.Retry(func() error {
		mysqlURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
			"test", "test", "localhost", resource.GetPort("3306/tcp"), "test")
		db, err = sql.Open("mysql", mysqlURI)
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	if err = migrateUp(migrationPath); err != nil {
		log.Fatalf("migration error: %s", err.Error())
	}

	return pool, resource
}

func closeDB(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func migrateUp(migrationPath string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath),
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	return m.Up()
}

func GetMysqlDB() *sql.DB {
	return db
}
