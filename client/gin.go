package client

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type GinEngine struct {
	addr   string
	debug  bool
	engine *gin.Engine
	logger zerolog.Logger
}

func NewGinEngine(addr string, port uint, debug bool, l zerolog.Logger) *GinEngine {
	l = l.With().Str("layer", "infrastructure").Str("component", "gin").Logger()

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
