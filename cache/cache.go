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
	Incr(key string) error
	Decr(key string) error
	GetKeys(pattern string) ([]string, error)
}
