package health

import (
	"github.com/bloock/go-kit/client"
	_ "github.com/go-sql-driver/mysql"
)

type HealthMysql struct {
	client      *client.MysqlClient
	description string
	version     string
}

func NewHealthMysql(client *client.MysqlClient, description string) HealthMysql {
	return HealthMysql{
		client:      client,
		description: description,
		version:     DepVersion("go-sql-driver/mysql"),
	}
}

func (h HealthMysql) HealthCheck() ExternalServiceDetails {
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
