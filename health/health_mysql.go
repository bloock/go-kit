package health

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type HealthMysql struct {
	db          *sql.DB
	description string
	version     string
}

func NewHealthMysql(db *sql.DB, description, version string) HealthMysql {
	return HealthMysql{
		db:          db,
		description: description,
		version:     version,
	}
}

func (h HealthMysql) HealthCheck() ExternalServiceDetails {
	s := "pass"
	var e string
	if err := h.db.Ping(); err != nil {
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
