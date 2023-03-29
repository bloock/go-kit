package cache

import (
	"context"
	"github.com/bloock/go-kit/domain"
)

type CacheUsageRepository interface {
	Save(ctx context.Context, usage domain.CacheUsage) error

	GetValueByKey(ctx context.Context, key string) (domain.CacheUsage, error)
	FindValueByKey(ctx context.Context, key string) (domain.CacheUsage, error)

	Update(ctx context.Context, usage domain.CacheUsage) error
}
