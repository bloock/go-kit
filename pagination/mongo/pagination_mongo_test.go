package pagination_mongo

import (
	"context"
	"github.com/bloock/go-kit/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationMongo(t *testing.T) {

	_mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer _mt.Close()

	_mt.Run("mongo pagination should parse ok", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"metadata", bson.A{
				bson.D{{"total", int64(10)}, {"page", int64(3)}},
			}},
			{"data", bson.A{
				bson.D{{"hello", "world"}},
			}},
		}))

		pq := pagination.PaginationQuery{
			Page:    3,
			PerPage: 15,
		}

		var res []struct {
			Hello string `bson:"hello"`
		}
		p, err := FindWithPagination(context.Background(), mt.Coll, nil, nil, pq, &res)

		assert.NoError(t, err)
		assert.Equal(t, "world", res[0].Hello)
		assert.Equal(t, int64(10), p.Meta.Total)
		assert.Equal(t, int64(3), p.Meta.CurrentPage)
		assert.Equal(t, int64(15), p.Meta.PerPage)
	})

	_mt.Run("mongo pagination with empty result should parse ok", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"metadata", bson.A{
				bson.D{{"total", int64(10)}, {"page", int64(3)}},
			}},
			{"data", bson.A{}},
		}))

		pq := pagination.PaginationQuery{
			Page:    3,
			PerPage: 15,
		}

		var res []struct {
			Hello string `bson:"hello"`
		}
		p, err := FindWithPagination(context.Background(), mt.Coll, nil, nil, pq, &res)

		assert.NoError(t, err)
		assert.Equal(t, 0, len(res))
		assert.Equal(t, int64(10), p.Meta.Total)
		assert.Equal(t, int64(3), p.Meta.CurrentPage)
		assert.Equal(t, int64(15), p.Meta.PerPage)
	})
}
