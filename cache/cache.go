package cache

import (
	"time"
)

//go:generate mockgen -source=cache/cache.go -destination mocks/cache/mock_cache.go
type Cache interface {
	TTL() time.Duration
	Set(key string, data []byte, expiration time.Duration) error
	SetInt(key string, data int) error
	Incr(key string) (int64, error)
	IncrBy(key string, quantity int) (int64, error)
	Decr(key string) (int64, error)
	Get(key string) ([]byte, error)
	GetInt(key string) (int, error)
	Del(key string) error
	DeleteKeys(keys []string) error
	ZAdd(key string, score float64, value []byte) error
	ZRem(key string, value []byte) error
	ZRangeByScore(key string, now string) ([]string, error)
	ZCount(key string, now string) (int64, error)
	MSet(keys []string, values []int32) error
	MGet(keys []string) ([]interface{}, error)
}
