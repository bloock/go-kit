package pagination_mysql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bloock/go-kit/pagination"
	"github.com/stretchr/testify/assert"
)

func TestPaginationMysql_GetPagination(t *testing.T) {

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	t.Run("Get pagination should return ok", func(t *testing.T) {
		sqlMock.ExpectQuery(
			`SELECT count(*) as total FROM table`).
			WillReturnRows(sqlmock.NewRows([]string{"total"}))

		_, err := GetPagination(context.Background(), db, "table", pagination.PaginationQuery{Page: 1, PerPage: 10})

		assert.NoError(t, sqlMock.ExpectationsWereMet())
		assert.NoError(t, err)
	})
}
