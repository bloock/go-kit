package containers

import (
	"context"
)

type ContainerOptions struct {
	Environment  []string
	Expiration   int
	Tag          string
	PortBindings map[string]string
	Host         string
	ExposedPorts []string
}

type DockerContainer interface {
	StartContainer(ctx context.Context) error
}

type GenericDockerContainer struct {
	Name         string
	Repository   string
	Tag          string `default:"latest"`
	Environment  []string
	Host         string `default:"localhost"`
	PortBindings map[string]string
	ExposedPorts []string
}
