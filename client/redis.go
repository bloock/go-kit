package client

import (
	"crypto/tls"
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

func NewRedis(host, port, password string, db int, enableTls bool, duration time.Duration, l zerolog.Logger) (*Redis, error) {
	l = l.With().Str("layer", "infrastructure").Str("component", "redis").Logger()

	var tlsConfig *tls.Config
	if enableTls {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	options := &redis.Options{
		Addr:      fmt.Sprintf("%s:%s", host, port),
		Password:  password,
		DB:        db,
		TLSConfig: tlsConfig,
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

func (r Redis) SetInt(key string, data int) error {
	return r.client.Set(key, data, 0).Err()
}

func (r Redis) Incr(key string) (int64, error) {
	return r.client.Incr(key).Result()
}

func (r Redis) IncrBy(key string, quantity int) (int64, error) {
	return r.client.IncrBy(key, int64(quantity)).Result()
}

func (r Redis) Decr(key string) (int64, error) {
	return r.client.Decr(key).Result()
}

func (r Redis) Get(key string) ([]byte, error) {
	result, err := r.client.Get(key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return result, err
}

func (r Redis) GetInt(key string) (int, error) {
	result, err := r.client.Get(key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return result, err
}

func (r Redis) Del(key string) error {
	return r.client.Del(key).Err()
}

func (r Redis) ZAdd(key string, score float64, value []byte) error {
	arg := redis.Z{
		Score:  score,
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

func (r Redis) ZCount(key string, now string) (int64, error) {
	min := now
	max := "+inf"

	return r.client.ZCount(key, min, max).Result()
}

func (r Redis) MSet(keys []string, values []int32) error {
	var pairs []interface{}
	for i := range keys {
		pairs = append(pairs, keys[i], values[i])
	}
	if err := r.client.MSet(pairs...).Err(); err != nil {
		r.logger.Error().Err(err).Msg("")
		return err
	}
	return nil
}

func (r Redis) MGet(keys []string) ([]interface{}, error) {
	res, err := r.client.MGet(keys...).Result()
	if err != nil {
		r.logger.Error().Err(err).Msg("")
		return nil, err
	}
	return res, nil
}

func (r Redis) Client() *redis.Client {
	return r.client
}
