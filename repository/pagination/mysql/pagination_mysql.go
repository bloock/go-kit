package pagination_mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/bloock/go-kit/repository/pagination"
	"github.com/huandu/go-sqlbuilder"
)

type paginationResponse struct {
	Total int `db:"total"`
}

func GetPagination(ctx context.Context, db *sql.DB, sb *sqlbuilder.SelectBuilder, pq pagination.PaginationQuery) (pagination.Pagination, error) {
	paginationSQLStruct := sqlbuilder.NewStruct(new(paginationResponse))

	sb.Select("count(*) as total").OrderBy().Limit(-1).Offset(-1)

	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctxTimeout, query, args...)
	if err != nil {
		return pagination.Pagination{}, err
	}
	defer rows.Close()

	rows.Next()
	var pr paginationResponse
	rows.Scan(paginationSQLStruct.Addr(&pr)...)

	return pagination.NewPagination(pq.Page, pq.PerPage, pr.Total), nil
}
