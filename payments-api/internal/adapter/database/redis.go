package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jtonynet/go-payments-api/config"
	redis "github.com/redis/go-redis/v9"
)

/*
	font: https://github.com/redis/go-redis
*/

type RedisClient struct {
	ctx context.Context

	client     *redis.Client
	strategy   string
	expiration time.Duration
}

func NewRedisClient(cfg config.InMemoryDatabase) (*RedisClient, error) {
	strAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     strAddr,
		Password: cfg.Pass,
		DB:       cfg.DB,
		Protocol: cfg.Protocol,
	})

	Expiration := time.Duration(cfg.Expiration * int(time.Millisecond))

	return &RedisClient{
		ctx: context.Background(),

		client:     client,
		strategy:   cfg.Strategy,
		expiration: Expiration,
	}, nil
}

func (c *RedisClient) Readiness(_ context.Context) error {
	_, err := c.client.Ping(c.ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) GetStrategy(_ context.Context) (string, error) {
	return c.strategy, nil
}

func (c *RedisClient) Set(_ context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.client.Set(c.ctx, key, data, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) Get(_ context.Context, key string) (string, error) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return "", err
	}
	if val == "" {
		return "", errors.New("get data empty")
	}

	return val, nil
}

func (c *RedisClient) Delete(_ context.Context, key string) error {
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := c.client.Expire(c.ctx, key, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) GetDefaultExpiration(_ context.Context) (time.Duration, error) {
	return c.expiration, nil
}

func (c *RedisClient) GetClient(_ context.Context) (interface{}, error) {
	return c.client, nil
}
