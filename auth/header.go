package auth

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetBearerTokenHeader(ctx *gin.Context) string {
	auth := ctx.Request.Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	splitToken := strings.Split(auth, "Bearer")
	if len(splitToken) != 2 {
		return ""
	}

	return strings.TrimSpace(splitToken[1])
}

func GetApiKeyHeader(ctx *gin.Context) string {
	apiKey := ctx.Request.Header.Get("x-api-key")
	if apiKey == "" {
		return ""
	}

	return apiKey
}
