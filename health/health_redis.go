package health

import (
	"github.com/bloock/go-kit/client"
)

type HealthRedis struct {
	client      *client.Redis
	description string
	version     string
}

func NewHealthRedis(client *client.Redis, description string) HealthRedis {
	return HealthRedis{
		client:      client,
		description: description,
		version:     DepVersion("go-redis/redis"),
	}
}

func (h HealthRedis) HealthCheck() ExternalServiceDetails {
	s := "pass"
	var e string

	_, err := h.client.Client().Ping().Result()
	if err != nil {
		s = "error"
		e = err.Error()
	}

	return ExternalServiceDetails{
		Description: h.description,
		Version:     h.version,
		Status:      s,
		Error:       e,
	}
}
