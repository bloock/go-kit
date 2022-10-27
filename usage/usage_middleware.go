package usage

import (
	"errors"
	"fmt"
	"github.com/bloock/go-kit/auth"
	"github.com/bloock/go-kit/cache"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

const (
	LIMIT_SUFFIX   = "limit"
	CLIENT_ID      = "client_id"
	USAGE_QUANTITY = "usage_quantity"
	USAGE_DISABLE  = "usage_disable"
)

type UsageMiddleware struct {
	logger zerolog.Logger
	redis  cache.Cache
	service string
}

func NewUsageMiddleware(l zerolog.Logger, redis cache.Cache, service string) UsageMiddleware {
	return UsageMiddleware{
		logger: l,
		redis:  redis,
		service: service,
	}
}

func (u UsageMiddleware) CheckUsageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.Request.Header.Get(auth.CLIENT_ID_HEADER)
		if clientID == "" {
			err := errors.New("no clientID provided")
			u.logger.Error().Err(err).Msg("")
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		keyLimit := GenerateUsageLimitKey(clientID, u.service)
		key := GenerateUsageKey(clientID, u.service)
		c.Set(CLIENT_ID, clientID)

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
		key := GenerateUsageKey(clientID, u.service)

		if isDisallow || listErrors != nil {
			return
		}

		if isQuantity {
			quantity = q.(int)
		}

		_, err := u.redis.IncrBy(key, quantity)
		if err != nil {
			u.logger.Error().Err(err).Msg("")
			return
		}
	}
}

func GenerateUsageLimitKey(clientID string, service string) string {
	return fmt.Sprintf("%s:%s:%s", service, clientID, LIMIT_SUFFIX)
}

func GenerateUsageKey(clientID string, service string) string {
	return fmt.Sprintf("%s:%s", service, clientID)
}
