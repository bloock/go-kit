package usage

import (
	"fmt"
	"github.com/bloock/go-kit/auth"
	"github.com/bloock/go-kit/cache"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

const (
	LIMIT_SUFFIX = ":limit"
	CLIENT_ID    = "client_id"
	USAGE_QUANTITY = "usage_quantity"
	USAGE_DISABLE  = "usage_disable"
)

type UsageMiddleware struct {
	logger zerolog.Logger
	redis  cache.Cache
}

func NewUsageMiddleware(l zerolog.Logger, redis cache.Cache) UsageMiddleware {
	return UsageMiddleware{
		logger: l,
		redis:  redis,
	}
}

func (u UsageMiddleware) CheckUsageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := auth.GetBearerToken(c.Request.Header.Get(auth.AUTHORIZATION_HEADER))
		var claims auth.JWTClaims
		err := auth.DecodeJWTUnverified(jwtToken, &claims)
		if err != nil {
			u.logger.Error().Err(err).Msg("")
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("invalid token provided"))
			c.Abort()
			return
		}

		key := claims.ClientID
		keyLimit := GenerateUsageLimitKey(key)
		c.Set(CLIENT_ID, key)

		limit, err := u.redis.GetInt(keyLimit)
		if err != nil {
			u.logger.Error().Err(err).Msg("")
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(fmt.Sprintf("redis error: %s", err.Error())))
			c.Abort()
			return
		}
		if limit == -1 {
			return
		}

		consumed, err := u.redis.GetInt(key)
		if err != nil {
			u.logger.Error().Err(err).Msg("")
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(fmt.Sprintf("redis error: %s", err.Error())))
			c.Abort()
			return
		}

		if consumed >= limit {
			u.logger.Error().Err(err).Msg("")
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("limit consumed"))
			c.Abort()
			return
		}
	}
}

func (u UsageMiddleware) UpdateUsageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		listErrors := c.Errors.Errors()
		clientID := c.MustGet(CLIENT_ID).(string)
		_, isDisallow := c.Get(USAGE_DISABLE)
		q, isQuantity := c.Get(USAGE_QUANTITY)
		var quantity = 1

		if isDisallow || listErrors != nil {
			return
		}

		if isQuantity {
			quantity = q.(int)
		}

		_, err := u.redis.IncrBy(clientID, quantity)
		if err != nil {
			u.logger.Error().Err(err).Msg("")
			return
		}
	}
}

func GenerateUsageLimitKey(key string) string {
	return fmt.Sprintf("%s%s", key, LIMIT_SUFFIX)
}
