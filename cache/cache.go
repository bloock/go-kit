package cache

import (
	"time"
)

//go:generate mockgen -package=mocks -source=cache/cache.go -destination mocks/mock_cache.go
type Cache interface {
	TTL() time.Duration
	Set(key string, data []byte, expiration time.Duration) error
	Get(key string) ([]byte, error)
	Del(key string) error
}
