package pagination_mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/bloock/go-kit/pagination"
	"github.com/huandu/go-sqlbuilder"
)

type paginationResponse struct {
	Total int `db:"total"`
}

func GetPagination(ctx context.Context, db *sql.DB, table string, pq pagination.PaginationQuery) (pagination.Pagination, error) {
	paginationSQLStruct := sqlbuilder.NewStruct(new(paginationResponse))
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("count(*) as total").
		From(table)

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
