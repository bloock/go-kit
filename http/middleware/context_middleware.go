package middleware

import (
	"github.com/bloock/go-kit/context"
	"github.com/gin-gonic/gin"
)

func ContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		xClientID := ctx.Request.Header.Get(context.UserIDKey)
		xRequestID := ctx.Request.Header.Get(context.RequestIDKey)
		ctx.Set(context.RequestIDKey, xRequestID)
		ctx.Set(context.UserIDKey, xClientID)
		ctx.Next()
	}
}
