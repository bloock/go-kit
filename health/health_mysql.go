package health

import (
	"database/sql"
	"fmt"

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
	fmt.Printf("Ping db")
	if err := h.db.Ping(); err != nil {
		fmt.Printf("Error db %s", err.Error())
		s = "error"
		e = err.Error()
	}

	fmt.Printf("%+v", h.db.Ping())
	return ExternalServiceDetails{
		Description: h.description,
		Version:     h.version,
		Status:      s,
		Error:       e,
	}
}
