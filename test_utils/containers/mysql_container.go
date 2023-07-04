package containers

import (
	"context"
	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/observability"
	"github.com/bloock/go-kit/test_utils"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"strings"
	"time"
)

var StartDefaultMysqlContainer = func(ctx context.Context) error {
	return MysqlDockerContainer{
		ContainerOptions{
			Environment: []string{
				"MYSQL_ROOT_PASSWORD=admin",
				"MYSQL_DATABASE=bloock"},
			Expiration:   120,
			Tag:          "8.0.22",
			ExposedPorts: []string{"3306"},
			Host:         "localhost",
			PortBindings: map[string]string{"3306": "3306"},
		},
		MysqlConnOptions{
			User:     "root",
			Password: "admin",
			DBName:   "bloock",
		},
	}.StartContainer(ctx)
}

type MysqlConnOptions struct {
	User     string `default:"root"`
	Password string `default:"admin"`
	DBName   string `default:"mysql"`
}

type MysqlDockerContainer struct {
	ContainerOptions
	MysqlConnOptions
}

func (gc GenericDockerContainer) StartContainer() (*dockertest.Resource, *dockertest.Pool, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
		return nil, nil, err
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
		return nil, nil, err
	}
	portBindings := map[docker.Port][]docker.PortBinding{}

	for port, hostPort := range gc.PortBindings {
		for _, exposedPort := range gc.ExposedPorts {
			portBindings[docker.Port(port+"/tcp")] = []docker.PortBinding{
				{HostIP: exposedPort, HostPort: hostPort + "/tcp"},
			}
		}
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:         gc.Name,
		Repository:   gc.Repository,
		Tag:          gc.Tag,
		Env:          gc.Environment,
		PortBindings: portBindings,
		ExposedPorts: gc.ExposedPorts,
	}, func(config *docker.HostConfig) {
		config.AutoRemove = false
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, nil, err
	}
	return resource, pool, nil
}

func (mc MysqlDockerContainer) StartContainer(ctx context.Context) error {

	_ = mysql.SetLogger(test_utils.Logger{})
	id, _, _ := strings.Cut(uuid.NewString(), "-")
	resource, pool, err := GenericDockerContainer{
		Name:         "mysql-" + id,
		Repository:   "mysql",
		Tag:          mc.Tag,
		Environment:  mc.Environment,
		Host:         mc.Host,
		PortBindings: mc.PortBindings,
		ExposedPorts: mc.ExposedPorts,
	}.StartContainer()

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
		return err
	}
	if err := pool.Retry(func() error {
		mysqlClient, err := client.NewMysqlClient(ctx, mc.User, mc.Password, mc.Host,
			mc.ExposedPorts[0], mc.DBName, false, &client.SQLConnOpts{
				MaxConnLifeTime: 5 * time.Minute,
				MaxOpenConns:    10,
				MaxIdleConns:    10,
			}, observability.Logger{})
		if err != nil {
			return err
		}
		return mysqlClient.DB().Ping()
	}); err != nil {
		log.Fatalf("Could not connect to mysql: %s", err)
		return err
	}
	err = resource.Expire(uint(mc.Expiration))
	if err != nil {
		log.Fatalf("Could not set expire to resource: %s", err)
		return err
	}
	return nil
}
