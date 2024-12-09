package mongodb

import (
	"context"
	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/repository/pagination"
)

func GetPagination(ctx context.Context, client *client.MongoClient, collectionName string, filter interface{}, pq pagination.PaginationQuery) (pagination.Pagination, error) {
	ctxWithTimeout, cancel := client.ContextWithTimeout(ctx)
	defer cancel()

	collection := client.DB().Collection(collectionName)

	count, err := collection.CountDocuments(ctxWithTimeout, filter)
	if err != nil {
		return pagination.Pagination{}, err
	}

	return pagination.NewPagination(pq.Page, pq.PerPage, int(count)), nil
}
