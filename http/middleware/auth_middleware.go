package middleware

import (
	"errors"
	"github.com/bloock/go-kit/auth"
	httpError "github.com/bloock/go-kit/errors"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ability auth.Ability) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header.Get(auth.AUTHORIZATION_HEADER)
		if authorizationHeader == "" {
			_ = c.Error(httpError.ErrUnauthorized(errors.New("no authorization header found")))
			c.Abort()
			return
		}

		jwtToken := auth.GetBearerToken(authorizationHeader)
		if jwtToken == "" {
			_ = c.Error(httpError.ErrUnauthorized(errors.New("invalid token provided")))
			c.Abort()
			return
		}

		var claims auth.JWTClaims
		err := auth.DecodeJWTUnverified(jwtToken, &claims)
		if err != nil {
			_ = c.Error(httpError.ErrUnauthorized(errors.New("invalid token content provided")))
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

		_ = c.Error(httpError.ErrForbidden(errors.New("action forbidden")))
		c.Abort()
	}
}
