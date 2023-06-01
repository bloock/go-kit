package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/bloock/go-kit/auth"
	"github.com/bloock/go-kit/cache"
	"github.com/bloock/go-kit/domain"
	httpError "github.com/bloock/go-kit/errors"
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
	logger          observability.Logger
	redis           cache.Cache
	usageRepository cache.CacheUsageRepository
	service         string
}

func NewUsageMiddleware(l observability.Logger, redis cache.Cache, cu cache.CacheUsageRepository, service string) UsageMiddleware {
	l.UpdateLogger(l.With().Caller().Str("component", "usage-middleware").Logger())

	return UsageMiddleware{
		logger:          l,
		redis:           redis,
		usageRepository: cu,
		service:         service,
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
				_ = c.Error(httpError.ErrUnauthorized(errors.New("invalid token provided")))
				u.logger.Info(c).Err(err).Msg("")
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
			_ = c.Error(err)
			u.logger.Info(c).Err(err).Msg("")
			c.Abort()
			return
		}
		if limit == -2 {
			limit, err = u.cacheMiss(c, keyLimit)
			if err != nil {
				_ = c.Error(err)
				c.Abort()
				return
			}
		}
		if limit == -1 {
			return
		}

		consumed, err := u.redis.GetInt(c, key)
		if err != nil {
			_ = c.Error(err)
			u.logger.Info(c).Err(err).Msg("")
			c.Abort()
			return
		}
		if consumed == -2 {
			consumed, err = u.cacheMiss(c, key)
			if err != nil {
				_ = c.Error(err)
				c.Abort()
				return
			}
		}

		if consumed >= limit {
			_ = c.Error(httpError.ErrForbidden(errors.New("limit consumed")))
			u.logger.Info(c).Err(err).Msg("")
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

		err := u.incrementByKey(c, key, quantity)
		if err != nil {
			u.logger.Info(c).Err(err).Msg("")
			return
		}
	}
}

func (u UsageMiddleware) cacheMiss(ctx context.Context, key string) (int, error) {
	cacheUsage, err := u.usageRepository.GetValueByKey(ctx, key)
	if err != nil {
		return 0, err
	}

	if err = u.redis.SetInt(ctx, cacheUsage.Key(), cacheUsage.Value()); err != nil {
		return 0, err
	}

	return cacheUsage.Value(), nil
}

func (u UsageMiddleware) incrementByKey(ctx context.Context, key string, quantity int) error {
	cacheUsage, err := u.usageRepository.FindValueByKey(ctx, key)
	if err != nil {
		return err
	}
	newQuantity := cacheUsage.Value() + quantity

	if cacheUsage.Key() != "" {
		updateCacheUsage := domain.NewCacheUsage(cacheUsage.Key(), newQuantity)
		if err = u.usageRepository.Update(ctx, updateCacheUsage); err != nil {
			return err
		}
	} else {
		newCacheUsage := domain.NewCacheUsage(key, newQuantity)
		if err = u.usageRepository.Save(ctx, newCacheUsage); err != nil {
			return err
		}
	}

	return u.redis.Del(ctx, key)
}

func GenerateUsageLimitKey(clientID string, service string) string {
	return fmt.Sprintf("%s:%s:%s", service, clientID, LIMIT_SUFFIX)
}

func GenerateUsageKey(clientID string, service string) string {
	return fmt.Sprintf("%s:%s", service, clientID)
}
