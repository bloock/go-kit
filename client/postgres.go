package client

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bloock/go-kit/observability"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/lib/pq"
)

type PostgresSQLClient struct {
	db     *sql.DB
	logger observability.Logger
}

func NewPostgresClient(ctx context.Context, user, pass, host, port, dbName string, ssl bool, connOpts *SQLConnOpts, l observability.Logger) (*PostgresSQLClient, error) {
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	l.UpdateLogger(l.With().Str("layer", "infrastructure").Str("component", "postgres").Logger())
	sslMode := "disable"
	if ssl {
		sslMode = "require"
	}

	postgresURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, dbName, sslMode)
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		l.Error(ctx).Msgf("error opening postgres on uri %s: %s", postgresURI, err.Error())
		return nil, err
	}

	db.SetConnMaxLifetime(connOpts.MaxConnLifeTime)
	db.SetMaxOpenConns(connOpts.MaxOpenConns)
	db.SetMaxIdleConns(connOpts.MaxIdleConns)

	return &PostgresSQLClient{
		db:     db,
		logger: l,
	}, nil
}

func (c PostgresSQLClient) DB() *sql.DB {
	return c.db
}

func (c PostgresSQLClient) MigrateUp(ctx context.Context, path string) error {
	driver, err := postgres.WithInstance(c.db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		"postgres",
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
