package cache

import (
	"context"
	"time"
)

//go:generate mockgen -source=cache/cache.go -destination mocks/cache/mock_cache.go
type Cache interface {
	TTL() time.Duration
	Set(ctx context.Context, key string, data []byte, expiration time.Duration) error
	SetInt(ctx context.Context, key string, data int) error
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, quantity int) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	Get(ctx context.Context, key string) ([]byte, error)
	GetInt(ctx context.Context, key string) (int, error)
	Del(ctx context.Context, key string) error
	DeleteKeys(ctx context.Context, keys []string) error
	ZAdd(ctx context.Context, key string, score float64, value []byte) error
	ZRem(ctx context.Context, key string, value []byte) error
	ZRangeByScore(ctx context.Context, key string, now string) ([]string, error)
	ZCount(ctx context.Context, key string, now string) (int64, error)
	MSet(ctx context.Context, keys []string, values []int32) error
	MGet(ctx context.Context, keys []string) ([]interface{}, error)
}
