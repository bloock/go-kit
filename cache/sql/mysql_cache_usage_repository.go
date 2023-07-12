package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bloock/go-kit/domain"
	"github.com/bloock/go-kit/errors"
	"github.com/bloock/go-kit/observability"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

type MysqlCacheUsageRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
	logger    observability.Logger
	service   string
}

func NewMysqlCacheUsageRepository(db *sql.DB, dbTimeout time.Duration, l observability.Logger, service string) *MysqlCacheUsageRepository {
	l.UpdateLogger(l.With().Caller().Str("component", "cache-usage-mysql").Logger())

	return &MysqlCacheUsageRepository{
		db:        db,
		dbTimeout: dbTimeout,
		logger:    l,
		service:   service,
	}
}

const SqlCacheUsageTable = "cache_usage"

func (c MysqlCacheUsageRepository) Save(ctx context.Context, usage domain.CacheUsage) error {
	span, ctx := observability.NewSpan(ctx, fmt.Sprintf("%s.cache-usage-repository.save", c.service))
	defer span.Finish()

	cacheSqlStruct := sqlbuilder.NewStruct(new(SqlCacheUsage))
	query, args := cacheSqlStruct.InsertInto(SqlCacheUsageTable, MapToSqlCacheUsage(usage)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, c.dbTimeout)
	defer cancel()

	_, err := c.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		err = errors.WrapSqlRepositoryError(err)
		c.logger.Info(ctx).Err(err).Msg("")
	}
	return err
}

func (c MysqlCacheUsageRepository) GetValueByKey(ctx context.Context, key string) (domain.CacheUsage, error) {
	span, ctx := observability.NewSpan(ctx, fmt.Sprintf("%s.cache-usage-repository.get-value-by-key", c.service))
	defer span.Finish()

	cacheSQLStruct := sqlbuilder.NewStruct(new(SqlCacheUsage))
	sb := cacheSQLStruct.SelectFrom(SqlCacheUsageTable)
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

func (c MysqlCacheUsageRepository) FindValueByKey(ctx context.Context, key string) (domain.CacheUsage, error) {
	span, ctx := observability.NewSpan(ctx, fmt.Sprintf("%s.cache-usage-repository.find-value-by-key", c.service))
	defer span.Finish()

	cacheSQLStruct := sqlbuilder.NewStruct(new(SqlCacheUsage))
	sb := cacheSQLStruct.SelectFrom(SqlCacheUsageTable)
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

func (c MysqlCacheUsageRepository) Update(ctx context.Context, usage domain.CacheUsage) error {
	span, ctx := observability.NewSpan(ctx, fmt.Sprintf("%s.cache-usage-repository.update", c.service))
	defer span.Finish()

	ub := sqlbuilder.NewUpdateBuilder()
	ub.Update(SqlCacheUsageTable).Set(ub.Assign("value", usage.Value()),
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
