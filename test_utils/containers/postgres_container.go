package containers

import "strings"

import (
	"context"
	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/observability"
	"github.com/google/uuid"
	"log"
	"time"
)

var StartDefaultPostgresContainer = func(ctx context.Context) error {
	return PostgresDockerContainer{
		ContainerOptions{
			Environment: []string{
				"POSTGRES_USER=admin",
				"POSTGRES_PASSWORD=admin",
			},
			Expiration:   120,
			Tag:          "11",
			ExposedPorts: []string{"5432"},
			Host:         "localhost",
			PortBindings: map[string]string{"5431": "5432"},
		},
		PostgresConnOptions{
			User:     "admin",
			Password: "admin",
			DBName:   "bloock",
		},
	}.StartContainer(ctx)
}

type PostgresConnOptions struct {
	User     string `default:"admin"`
	Password string `default:"admin"`
	DBName   string `default:"bloock"`
}

type PostgresDockerContainer struct {
	ContainerOptions
	PostgresConnOptions
}

func (mc PostgresDockerContainer) StartContainer(ctx context.Context) error {

	id, _, _ := strings.Cut(uuid.NewString(), "-")
	resource, pool, err := GenericDockerContainer{
		Name:         "postgres-" + id,
		Repository:   "postgres",
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
		postgresClient, err := client.NewPostgresClient(ctx, mc.User, mc.Password, mc.Host,
			mc.ExposedPorts[0], mc.DBName, false, &client.SQLConnOpts{
				MaxConnLifeTime: 5 * time.Minute,
				MaxOpenConns:    10,
				MaxIdleConns:    10,
			}, observability.Logger{})
		if err != nil {
			return err
		}
		return postgresClient.DB().Ping()
	}); err != nil {
		log.Fatalf("Could not connect to postgres: %s", err)
		return err
	}
	err = resource.Expire(uint(mc.Expiration))
	if err != nil {
		log.Fatalf("Could not set expire to resource: %s", err)
		return err
	}
	return nil
}
