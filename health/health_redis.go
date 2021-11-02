package health

import (
	"database/sql"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type HealthRedis struct {
	client      *redis.Client
	description string
	version     string
}

func NewHealthRedis(db *sql.DB, description, version string) HealthMysql {
	return HealthMysql{
		db:          db,
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
