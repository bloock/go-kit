package containers

import (
	"context"
	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/observability"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

var StartDefaultRedisContainer = func(ctx context.Context) error {
	return RedisDockerContainer{
		ContainerOptions{
			Environment:  nil,
			Expiration:   120,
			Tag:          "latest",
			Host:         "localhost",
			ExposedPorts: []string{"6379"},
			PortBindings: map[string]string{"6379": "6379"},
		},
		RedisConnOptions{DB: 2},
	}.StartContainer(ctx)
}

type RedisConnOptions struct {
	DB int `default:"0"`
}

type RedisDockerContainer struct {
	ContainerOptions
	RedisConnOptions
}

func (rc RedisDockerContainer) StartContainer(ctx context.Context) error {
	id, _, _ := strings.Cut(uuid.NewString(), "-")
	resource, pool, err := GenericDockerContainer{
		Name:         "redis-" + id,
		Repository:   "redis",
		Tag:          rc.Tag,
		Environment:  rc.Environment,
		Host:         rc.Host,
		PortBindings: rc.PortBindings,
		ExposedPorts: rc.ExposedPorts,
	}.StartContainer()
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
		return err
	}

	if err := pool.Retry(func() error {
		_, err = client.NewRedis(ctx,
			rc.Host, rc.ExposedPorts[0], "", rc.DB, false, time.Duration(5)*time.Second, observability.Logger{},
		)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to redis: %s", err)
		return err
	}

	err = resource.Expire(uint(rc.Expiration))
	if err != nil {
		log.Fatalf("Could not set expire to resource: %s", err)
		return err
	}
	return nil
}
