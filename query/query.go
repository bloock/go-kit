package query

import "context"

type Bus interface {
	Dispatch(context.Context, Query) (interface{}, error)
	Register(Type, Handler)
}

//go:generate mockgen -package=mocks -source=query/query.go -destination mocks/mock_query.go
type Type string

type Query interface {
	Type() Type
}

type Handler interface {
	Handle(context.Context, Query) (interface{}, error)
}