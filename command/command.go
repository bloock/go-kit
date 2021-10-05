package command

import "context"

type Bus interface {
	Dispatch(context.Context, Command) error
	Register(Type, Handler)
}

//go:generate mockgen	-source=command/command.go -destination mocks/command/mock_command.go
type Type string

type Command interface {
	Type() Type
}

type Handler interface {
	Handle(context.Context, Command) error
}
