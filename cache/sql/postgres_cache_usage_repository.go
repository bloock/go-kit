package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/domain"
	"github.com/bloock/go-kit/errors"
	"github.com/bloock/go-kit/observability"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

const CACHE_USAGE_TABLE = "cache_usage"
const CACHE_SCHEMA = "cache"

type PostgresCacheUsageRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
	logger    observability.Logger
	service   string
}

func NewPostgresCacheUsageRepository(psql client.PostgresSQLClient, dbTimeout time.Duration, l observability.Logger, service string) *PostgresCacheUsageRepository {
	l.UpdateLogger(l.With().Caller().Str("component", "cache-usage-postgres").Logger())

	return &PostgresCacheUsageRepository{
		db:        psql.DB(),
		dbTimeout: dbTimeout,
		logger:    l,
		service:   service,
	}
}

func (c PostgresCacheUsageRepository) Save(ctx context.Context, usage domain.CacheUsage) error {
	cacheSqlStruct := sqlbuilder.NewStruct(new(SqlCacheUsage))
	query, args := cacheSqlStruct.InsertInto(fmt.Sprintf("%s.%s", CACHE_SCHEMA, CACHE_USAGE_TABLE), MapToSqlCacheUsage(usage)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()

	_, err := c.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		err = errors.WrapSqlRepositoryError(err)
		c.logger.Info(ctx).Err(err).Msg("")
	}
	return err
}

func (c PostgresCacheUsageRepository) GetValueByKey(ctx context.Context, key string) (domain.CacheUsage, error) {
	cacheSQLStruct := sqlbuilder.NewStruct(new(SqlCacheUsage))
	sb := cacheSQLStruct.SelectFrom(fmt.Sprintf("%s.%s", CACHE_SCHEMA, CACHE_USAGE_TABLE))
	sb = sb.Where(sb.Equal("_key", key))
	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()

	row := c.db.QueryRowContext(ctxTimeout, query, args...)

	var cc SqlCacheUsage
	err := row.Scan(cacheSQLStruct.Addr(&cc)...)
	if err != nil {
		err = errors.WrapSqlRepositoryError(err)
		c.logger.Info(ctx).Err(err).Msg(err.Error())
		return domain.CacheUsage{}, err
	}
	return MapToCacheUsage(cc), nil
}

func (c PostgresCacheUsageRepository) FindValueByKey(ctx context.Context, key string) (domain.CacheUsage, error) {
	cacheSQLStruct := sqlbuilder.NewStruct(new(SqlCacheUsage))
	sb := cacheSQLStruct.SelectFrom(fmt.Sprintf("%s.%s", CACHE_SCHEMA, CACHE_USAGE_TABLE))
	sb = sb.Where(sb.Equal("_key", key))
	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()

	rows, err := c.db.QueryContext(ctxTimeout, query, args...)
	if err != nil {
		err = errors.WrapSqlRepositoryError(err)
		c.logger.Info(ctx).Err(err).Msg("")
		return domain.CacheUsage{}, err
	}

	for rows.Next() {
		var cc SqlCacheUsage
		err = rows.Scan(cacheSQLStruct.Addr(&cc)...)
		if err != nil {
			err = errors.WrapSqlRepositoryError(err)
			c.logger.Info(ctx).Err(err).Msg(err.Error())
			return domain.CacheUsage{}, err
		}
		return MapToCacheUsage(cc), nil
	}

	return domain.CacheUsage{}, nil
}

func (c PostgresCacheUsageRepository) Update(ctx context.Context, usage domain.CacheUsage) error {
	ub := sqlbuilder.NewUpdateBuilder()
	ub.Update(fmt.Sprintf("%s.%s", CACHE_SCHEMA, CACHE_USAGE_TABLE)).Set(ub.Assign("value", usage.Value()),
		ub.Assign("updated_at", time.Now())).
		Where(ub.In("_key", usage.Key()))
	query, args := ub.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()

	_, err := c.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		err = errors.WrapSqlRepositoryError(err)
		c.logger.Info(ctx).Err(err).Msg("")
	}

	return err
}
