package client

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bloock/go-kit/observability"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MysqlClient struct {
	db     *sql.DB
	logger observability.Logger
}

func NewMysqlClient(ctx context.Context, user, pass, host, port, dbName string, tls bool, l observability.Logger) (*MysqlClient, error) {
	l.UpdateLogger(l.With().Str("layer", "infrastructure").Str("component", "mysql").Logger())

	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=%t&charset=utf8&parseTime=True", user, pass, host, port, dbName, tls)
	db, err := sql.Open("mysql", mysqlURI)
	if err != nil {
		l.Error(ctx).Msgf("error opening mysql on uri %s: %s", mysqlURI, err.Error())
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MysqlClient{
		db:     db,
		logger: l,
	}, nil
}

func (c MysqlClient) DB() *sql.DB {
	return c.db
}

func (c MysqlClient) MigrateUp(ctx context.Context, path string) error {
	driver, err := mysql.WithInstance(c.db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		"mysql",
		driver,
	)
	if err != nil {
		c.logger.Error(ctx).Msgf("migration error: %s", err.Error())
		return err
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			c.logger.Info(ctx).Msgf("no migration changes: %s", err.Error())
			return nil
		}

		c.logger.Error(ctx).Msgf("migration error: %s", err.Error())
		return err
	}

	return nil
}
