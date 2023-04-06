package containers

import (
	"context"
	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/observability"
	"github.com/google/uuid"
	"log"
	"strings"
)

var StartDefaultRabbitMQContainer = func(ctx context.Context) error {
	return RabbitMQDockerContainer{
		ContainerOptions: ContainerOptions{
			Environment:  []string{},
			Expiration:   120,
			Tag:          "3-management",
			ExposedPorts: []string{"5672"},
			Host:         "localhost",
			PortBindings: map[string]string{"5672": "5672"},
		},
		RabbitMQConnOptions: RabbitMQConnOptions{
			User:     "guest",
			Password: "guest",
		},
	}.StartContainer(ctx)
}

type RabbitMQConnOptions struct {
	User     string `default:"guest"`
	Password string `default:"guest"`
}
type RabbitMQDockerContainer struct {
	ContainerOptions
	RabbitMQConnOptions
}

func (rc RabbitMQDockerContainer) StartContainer(ctx context.Context) error {

	id, _, _ := strings.Cut(uuid.NewString(), "-")
	resource, pool, err := GenericDockerContainer{
		Name:         "rabbitmq-" + id,
		Repository:   "rabbitmq",
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

		_, err = client.NewAMQPClient(ctx, "bloock", rc.User, rc.Password, rc.Host, rc.ExposedPorts[0], "/", observability.Logger{})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to rabbitmq: %s", err)
	}

	err = resource.Expire(uint(rc.Expiration))
	if err != nil {
		log.Fatalf("Could not set expire to resource: %s", err)
		return err
	}
	return nil
}
