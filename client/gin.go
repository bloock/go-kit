package client

import (
	"fmt"

	"github.com/bloock/go-kit/observability"
	"github.com/gin-gonic/gin"
)

type GinEngine struct {
	addr   string
	debug  bool
	engine *gin.Engine
	logger observability.Logger
}

func NewGinEngine(addr string, port uint, debug bool, l observability.Logger) *GinEngine {
	l.UpdateLogger(l.With().Str("layer", "infrastructure").Str("component", "gin").Logger())

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = l.With().Str("level", "info").Logger()
	gin.DefaultErrorWriter = l.With().Str("level", "error").Logger()

	return &GinEngine{
		addr:   fmt.Sprintf("%s:%d", addr, port),
		debug:  debug,
		engine: gin.New(),
		logger: l,
	}
}

func (g GinEngine) Addr() string {
	return g.addr
}

func (g GinEngine) Engine() *gin.Engine {
	return g.engine
}

func (g GinEngine) Debug() bool {
	return g.debug
}
