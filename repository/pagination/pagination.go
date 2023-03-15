package pagination

import (
	"fmt"
	httpError "github.com/bloock/go-kit/errors"

	"github.com/gin-gonic/gin"
)

var FirstPage = 1

type PaginationQuery struct {
	Page    int `form:"page" json:"page"`
	PerPage int `form:"per_page" json:"per_page"`
}

func NewPaginationQuery(ctx *gin.Context) (PaginationQuery, error) {
	var pq PaginationQuery
	err := ctx.BindQuery(&pq)
	if err != nil {
		err = httpError.ErrInvalidBodyJSON(err)
		return PaginationQuery{}, err
	}

	if pq.Page < FirstPage {
		err = httpError.ErrInvalidBodyJSON(fmt.Errorf("page number should be bigger than %d", FirstPage))
		return PaginationQuery{}, err
	}

	if pq.PerPage < 1 {
		err = httpError.ErrInvalidBodyJSON(fmt.Errorf("per page value should be bigger than 1"))
		return PaginationQuery{}, err
	}

	return pq, nil
}

func (p PaginationQuery) Skip() int {
	skip := (p.Page - FirstPage) * p.PerPage
	return skip
}

type metaPagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	From        int `json:"from"`
	To          int `json:"to"`
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
}

type Pagination struct {
	Meta metaPagination `json:"meta"`
}

func NewPagination(currentPage, perPage, total int) Pagination {
	from := ((currentPage - FirstPage) * perPage) + 1

	var lastPage int
	if total == 0 {
		lastPage = 1
	} else {
		if total%perPage != 0 {
			lastPage = total/perPage + 1
		} else {
			lastPage = total / perPage
		}
	}

	var to int
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
