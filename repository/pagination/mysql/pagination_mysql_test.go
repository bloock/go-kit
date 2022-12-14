package pagination_mysql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bloock/go-kit/repository/pagination"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestPaginationMysql_GetPagination(t *testing.T) {

	table := "table"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	t.Run("Get pagination should return ok", func(t *testing.T) {
		sb := sqlbuilder.NewSelectBuilder()
		sb.Select("*").
			From(table).
			Where(sb.Equal("arg", true)).
			Offset(1).
			Limit(10)

		query, _ := sb.Build()

		assert.Equal(t, "SELECT * FROM table WHERE arg = ? LIMIT 10 OFFSET 1", query)

		sqlMock.ExpectQuery(
			`SELECT count(*) as total FROM table WHERE arg = ?`).
			WillReturnRows(sqlmock.NewRows([]string{"total"}))

		_, err := GetPagination(context.Background(), db, sb, pagination.PaginationQuery{Page: 1, PerPage: 10})

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.NoError(t, err)
	})
}
