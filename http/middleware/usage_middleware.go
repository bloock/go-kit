package middleware

import (
	"fmt"
	"net/http"

	"github.com/bloock/go-kit/auth"
	"github.com/bloock/go-kit/cache"
	"github.com/bloock/go-kit/observability"
	"github.com/gin-gonic/gin"
)

const (
	LIMIT_SUFFIX   = "limit"
	CLIENT_ID      = "client_id"
	USAGE_QUANTITY = "usage_quantity"
	USAGE_DISABLE  = "usage_disable"

	CoreService            = "core"
	NodeService            = "node"
	StorageService         = "storage"
	TransferService        = "transfer"
	KeysTransactionService = "keys_transaction"
)

type UsageMiddleware struct {
	logger  observability.Logger
	redis   cache.Cache
	service string
}

func NewUsageMiddleware(l observability.Logger, redis cache.Cache, service string) UsageMiddleware {
	return UsageMiddleware{
		logger:  l,
		redis:   redis,
		service: service,
	}
}

func (u UsageMiddleware) CheckUsageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.Request.Header.Get(auth.CLIENT_ID_HEADER)
		if clientID == "" {
			jwtToken := auth.GetBearerTokenHeader(c)
			var claims auth.JWTClaims
			err := auth.DecodeJWTUnverified(jwtToken, &claims)
			if err != nil {
				u.logger.Error(c).Err(err).Msg("")
				c.Writer.WriteHeader(http.StatusUnauthorized)
				c.Writer.Write([]byte("invalid token provided"))
				c.Abort()
				return
			}
			clientID = claims.ClientID
		}

		keyLimit := GenerateUsageLimitKey(clientID, u.service)
		key := GenerateUsageKey(clientID, u.service)
		c.Set(CLIENT_ID, clientID)

		limit, err := u.redis.GetInt(c, keyLimit)
		if err != nil {
			u.logger.Error(c).Err(err).Msg("")
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(fmt.Sprintf("redis error: %s", err.Error())))
			c.Abort()
			return
		}
		if limit == -1 {
			return
		}

		consumed, err := u.redis.GetInt(c, key)
		if err != nil {
			u.logger.Error(c).Err(err).Msg("")
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(fmt.Sprintf("redis error: %s", err.Error())))
			c.Abort()
			return
		}

		if consumed >= limit {
			u.logger.Error(c).Err(err).Msg("")
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

		_, err := u.redis.IncrBy(c, key, quantity)
		if err != nil {
			u.logger.Error(c).Err(err).Msg("")
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
