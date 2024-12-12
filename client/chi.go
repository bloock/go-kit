package client

import (
	"fmt"
	"github.com/bloock/go-kit/observability"
	"github.com/go-chi/chi/v5"
)

type ChiRouter struct {
	addr   string
	debug  bool
	router *chi.Mux
	logger observability.Logger
}

func NewChiRouter(addr string, port uint, debug bool, l observability.Logger) *ChiRouter {
	l.UpdateLogger(l.With().Str("layer", "infrastructure").Str("component", "chi").Logger())

	return &ChiRouter{
		addr:   fmt.Sprintf("%s:%d", addr, port),
		debug:  debug,
		router: chi.NewRouter(),
		logger: l,
	}
}

func (c ChiRouter) Addr() string {
	return c.addr
}

func (c ChiRouter) Router() *chi.Mux {
	return c.router
}

func (c ChiRouter) Debug() bool {
	return c.debug
}
