package health

import (
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type HealthRedis struct {
	client      *redis.Client
	description string
	version     string
}

func NewHealthRedis(client *redis.Client, description, version string) HealthRedis {
	return HealthRedis{
		client:      client,
		description: description,
		version:     version,
	}
}

func (h HealthRedis) HealthCheck() ExternalServiceDetails {
	s := "pass"
	var e string

	_, err := h.client.Ping().Result()
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
