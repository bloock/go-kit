package domain

type CacheUsage struct {
	key   string
	value int
}

func NewCacheUsage(key string, value int) CacheUsage {
	return CacheUsage{
		key:   key,
		value: value,
	}
}

func (c CacheUsage) Key() string {
	return c.key
}

func (c CacheUsage) Value() int {
	return c.value
}
