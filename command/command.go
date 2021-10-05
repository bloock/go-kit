package command

import "context"

type Bus interface {
	Dispatch(context.Context, Command) error
	Register(Type, Handler)
}

//go:generate mockgen -package=mocks -source=command/command.go -destination mocks/mock_command.go
type Type string

type Command interface {
	Type() Type
}

type Handler interface {
	Handle(context.Context, Command) error
}
