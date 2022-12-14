package middleware

import (
	"net/http"

	"github.com/bloock/go-kit/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ability auth.Ability) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header.Get(auth.AUTHORIZATION_HEADER)
		if authorizationHeader == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("no authorization header found"))
			c.Abort()
			return
		}

		jwtToken := auth.GetBearerToken(authorizationHeader)
		if jwtToken == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("invalid token provided"))
			c.Abort()
			return
		}

		var claims auth.JWTClaims
		err := auth.DecodeJWTUnverified(jwtToken, &claims)
		if err != nil {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("invalid token content provided"))
			c.Abort()
			return
		}

		if allowedActions := claims.Scopes[ability.Resource()]; allowedActions != nil {
			for _, a := range allowedActions {
				if a == ability.Action() {
					c.Next()
					return
				}
			}
		}

		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte("action forbbiden"))
		c.Abort()
	}
}
