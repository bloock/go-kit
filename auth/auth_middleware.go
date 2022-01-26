package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AUTHORIZATION_HEADER = "Authorization"
	BEARER_PREFIX        = "Bearer"
)

type Auth struct {
	jwtSecret string
}

func NewAuth(jwtSecret string) Auth {
	return Auth{
		jwtSecret: jwtSecret,
	}
}

func (a *Auth) Middleware(ability Ability) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header.Get(AUTHORIZATION_HEADER)
		if authorizationHeader == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("no authorization header found"))
			c.Abort()
			return
		}

		var claims JWTClaims
		err := DecodeJWT(authorizationHeader, a.jwtSecret, &claims)
		if err != nil {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("invalid token provided"))
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
