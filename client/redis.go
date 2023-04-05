package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/bloock/go-kit/observability"
	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
	ttl    time.Duration
	logger observability.Logger
}

func NewRedis(ctx context.Context, host, port, password string, db int, enableTls bool, duration time.Duration, l observability.Logger) (*Redis, error) {
	l.UpdateLogger(l.With().Str("layer", "infrastructure").Str("component", "redis").Logger())

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
		l.Error(ctx).Msg("couldn't connect to redis server")
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

func (r Redis) Set(ctx context.Context, key string, data []byte, expiration time.Duration) error {
	if err := r.client.Set(key, data, expiration).Err(); err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return err
	}
	return nil
}

func (r Redis) SetInt(ctx context.Context, key string, data int) error {
	if err := r.client.Set(key, data, 0).Err(); err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return err
	}
	return nil
}

func (r Redis) Incr(ctx context.Context, key string) (int64, error) {
	res, err := r.client.Incr(key).Result()
	if err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return res, err
	}
	return res, nil
}

func (r Redis) IncrBy(ctx context.Context, key string, quantity int) (int64, error) {
	res, err := r.client.IncrBy(key, int64(quantity)).Result()
	if err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return res, err
	}
	return res, nil
}

func (r Redis) Decr(ctx context.Context, key string) (int64, error) {
	res, err := r.client.Decr(key).Result()
	if err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return res, err
	}
	return res, nil
}

func (r Redis) DecrBy(ctx context.Context, key string, quantity int64) (int64, error) {
	res, err := r.client.DecrBy(key, quantity).Result()
	if err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return res, err
	}
	return res, nil
}

func (r Redis) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := r.client.Get(key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return nil, err
	}
	return result, nil
}

func (r Redis) GetInt(ctx context.Context, key string) (int, error) {
	result, err := r.client.Get(key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return 0, err
	}
	return result, nil
}

func (r Redis) Del(ctx context.Context, key string) error {
	if err := r.client.Del(key).Err(); err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return err
	}
	return nil
}

func (r Redis) DeleteKeys(ctx context.Context, keys []string) error {
	if err := r.client.Del(keys...).Err(); err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return err
	}
	return nil
}

func (r Redis) ZAdd(ctx context.Context, key string, score float64, value []byte) error {
	arg := redis.Z{
		Score:  score,
		Member: value,
	}
	if err := r.client.ZAdd(key, arg).Err(); err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return err
	}
	return nil
}

func (r Redis) ZRem(ctx context.Context, key string, value []byte) error {
	if err := r.client.ZRem(key, value).Err(); err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return err
	}
	return nil
}

func (r Redis) ZRangeByScore(ctx context.Context, key string, now string) ([]string, error) {
	arg := redis.ZRangeBy{
		Min: now,
		Max: "+inf",
	}
	res, err := r.client.ZRangeByScore(key, arg).Result()
	if err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return res, err
	}
	return res, nil
}

func (r Redis) ZCount(ctx context.Context, key string, now string) (int64, error) {
	min := now
	max := "+inf"

	res, err := r.client.ZCount(key, min, max).Result()
	if err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return res, err
	}
	return res, nil
}

func (r Redis) MSet(ctx context.Context, keys []string, values []int32) error {
	var pairs []interface{}
	for i := range keys {
		pairs = append(pairs, keys[i], values[i])
	}
	if err := r.client.MSet(pairs...).Err(); err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return err
	}
	return nil
}

func (r Redis) MGet(ctx context.Context, keys []string) ([]interface{}, error) {
	res, err := r.client.MGet(keys...).Result()
	if err != nil {
		r.logger.Info(ctx).Msg(err.Error()).Msg("")
		return nil, err
	}
	return res, nil
}

func (r Redis) Client() *redis.Client {
	return r.client
}
