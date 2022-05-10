package client

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog"
)

type Redis struct {
	client *redis.Client
	ttl    time.Duration
	logger zerolog.Logger
}

func NewRedis(host, port, password string, db int, duration time.Duration, l zerolog.Logger) (*Redis, error) {
	l = l.With().Str("layer", "infrastructure").Str("component", "redis").Logger()

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	}
	client := redis.NewClient(options)
	_, err := client.Ping().Result()
	if err != nil {
		l.Error().Msg("couldn't connect to redis server")
		return nil, err
	}

	return &Redis{
		client: client,
		ttl:    duration,
		logger: l,
	}, nil
}

func (r Redis) TTL() time.Duration {
	return r.ttl
}

func (r Redis) Set(key string, data []byte, expiration time.Duration) error {
	return r.client.Set(key, data, expiration).Err()
}

func (r Redis) Get(key string) ([]byte, error) {
	result, err := r.client.Get(key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return result, err
}

func (r Redis) Del(key string) error {
	return r.client.Del(key).Err()
}

func (r Redis) Client() *redis.Client {
	return r.client
}

func (r Redis) ZAdd(key string, score float64, value []byte) error {
	arg := redis.Z{
		Score: score,
		Member: value,
	}
	return r.client.ZAdd(key, arg).Err()
}

func (r Redis) ZRem(key string, value []byte) error {
	return r.client.ZRem(key, value).Err()
}

func (r Redis) ZRangeByScore(key string, now string) ([]string, error) {
	arg := redis.ZRangeBy{
		Min: now,
		Max: "+inf",
	}
	return r.client.ZRangeByScore(key, arg).Result()
}