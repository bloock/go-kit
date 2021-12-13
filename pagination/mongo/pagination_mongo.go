package pagination_mongo

import (
	"context"
	"fmt"
	"github.com/bloock/go-kit/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type paginationResponse struct {
	Metadata []struct {
		Total int64 `bson:"total"`
		Page  int64 `bson:"page"`
	} `bson:"metadata"`
	Data []bson.M `bson:"data"`
}

func FindWithPagination(coll *mongo.Collection, ctx context.Context, filter bson.D, pq pagination.PaginationQuery, res interface{}) (pagination.Pagination, error) {
	match := bson.D{
		{"$match", filter},
	}
	facet := bson.D{
		{"$facet", bson.D{
			{"metadata", bson.A{
				bson.D{{"$count", "total"}},
				bson.D{{"$addFields", bson.D{{"page", pq.Page}}}},
			}},
			{"data", bson.A{
				bson.D{{"$skip", pq.Skip()}},
				bson.D{{"$limit", pq.PerPage}},
			}},
		}},
	}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{match, facet})
	if err != nil {
		return pagination.Pagination{}, err
	}
	defer cursor.Close(ctx)

	var response paginationResponse
	if next := cursor.Next(ctx); !next {
		return pagination.Pagination{}, err
	}

	if err = cursor.Decode(&response); err != nil {
		return pagination.Pagination{}, err
	}

	var p pagination.Pagination
	if len(response.Metadata) < 1 {
		p = pagination.NewPagination(1, pq.PerPage, 0)
	} else {
		p = pagination.NewPagination(response.Metadata[0].Page, pq.PerPage, response.Metadata[0].Total)
	}

	if err = decodeResult(response.Data, res); err != nil {
		return pagination.Pagination{}, err
	}

	return p, nil
}

func decodeResult(val []bson.M, res interface{}) error {
	resultsVal := reflect.ValueOf(res)
	if resultsVal.Kind() != reflect.Ptr {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a %s", resultsVal.Kind())
	}

	sliceVal := resultsVal.Elem()
	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a pointer to %s", sliceVal.Kind())
	}

	elementType := sliceVal.Type().Elem()

	index := 0

	for _, item := range val {
		bsonBytes, err := bson.Marshal(item)
		if err != nil {
			return err
		}

		if sliceVal.Len() == index {
			// slice is full
			newElem := reflect.New(elementType)
			sliceVal = reflect.Append(sliceVal, newElem.Elem())
			sliceVal = sliceVal.Slice(0, sliceVal.Cap())
		}

		currElem := sliceVal.Index(index).Addr().Interface()
		if err = bson.Unmarshal(bsonBytes, currElem); err != nil {
			return err
		}

		index++
	}

	resultsVal.Elem().Set(sliceVal.Slice(0, index))
	return nil
}
