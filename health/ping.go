package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:route GET /ping infrastructure ping
//
//	Schemes: http, https
//
func PingHandler() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	}
}
