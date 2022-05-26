package cache

import (
	"time"
)

//go:generate mockgen -source=cache/cache.go -destination mocks/cache/mock_cache.go
type Cache interface {
	TTL() time.Duration
	Set(key string, data []byte, expiration time.Duration) error
	Get(key string) ([]byte, error)
	Del(key string) error
	ZAdd(key string, score float64, value []byte) error
	ZRem(key string, value []byte) error
	ZRangeByScore(key string, now string) ([]string, error)
	ZCount(key string, now string) (int64, error)
}
