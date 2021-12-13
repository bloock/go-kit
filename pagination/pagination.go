package pagination

import (
	"github.com/gin-gonic/gin"
)

var FirstPage = int64(1)

type PaginationQuery struct {
	Page    int64 `form:"page" json:"page"`
	PerPage int64 `form:"per_page" json:"per_page"`
}

func NewPaginationQuery(ctx *gin.Context) (PaginationQuery, error) {
	var pq PaginationQuery
	err := ctx.BindQuery(&pq)
	if err != nil {
		return PaginationQuery{}, err
	}

	return pq, nil
}

func (p PaginationQuery) Skip() *int64 {
	skip := (p.Page - FirstPage) * p.PerPage
	return &skip
}

type metaPagination struct {
	CurrentPage int64 `json:"current_page"`
	PerPage     int64 `json:"per_page"`
	From        int64 `json:"from"`
	To          int64 `json:"to"`
	Total       int64 `json:"total"`
	LastPage    int64 `json:"last_page"`
}

type Pagination struct {
	Meta metaPagination `json:"meta"`
}

func NewPagination(currentPage, perPage, total int64) Pagination {
	from := ((currentPage - FirstPage) * perPage) + 1

	var lastPage int64
	if total == 0 {
		lastPage = 1
	} else {
		if total%perPage != 0 {
			lastPage = total/perPage + 1
		} else {
			lastPage = total / perPage
		}
	}

	var to int64
	if total == 0 {
		to = 1
	} else {
		if currentPage == lastPage {
			if perPage <= total {
				to = total
			} else {
				to = (currentPage-FirstPage)*perPage + (total % perPage)
			}
		} else {
			to = currentPage * perPage
		}
	}

	return Pagination{
		Meta: metaPagination{
			CurrentPage: currentPage,
			PerPage:     perPage,
			From:        from,
			To:          to,
			Total:       total,
			LastPage:    lastPage,
		},
	}
}
