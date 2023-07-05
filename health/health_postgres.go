package health

import (
	"github.com/bloock/go-kit/client"
	_ "github.com/go-sql-driver/mysql"
)

type HealthPostgres struct {
	client      *client.PostgresSQLClient
	description string
	version     string
}

func NewHealthPostgres(client *client.PostgresSQLClient, description string) HealthPostgres {
	return HealthPostgres{
		client:      client,
		description: description,
		version:     DepVersion("go-sql-driver/postgres"),
	}
}

func (h HealthPostgres) HealthCheck() ExternalServiceDetails {
	s := "pass"
	var e string
	if err := h.client.DB().Ping(); err != nil {
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
