package pubSub

import (
	"context"
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	client *redis.Client
	pubsub *redis.PubSub

	strategy string
}

func NewRedisPubSub(cfg config.PubSub) (*RedisPubSub, error) {
	strAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     strAddr,
		Password: cfg.Pass,
		DB:       cfg.DB,
		Protocol: cfg.Protocol,
	})

	return &RedisPubSub{
		client:   client,
		strategy: cfg.Strategy,
	}, nil
}

func (r *RedisPubSub) Subscribe(_ context.Context, key string) (<-chan string, error) {
	keyspaceChannel := fmt.Sprintf("__keyevent@%d__:expired", r.client.Options().DB)
	bufferSize := 500

	r.pubsub = r.client.Subscribe(context.Background(), keyspaceChannel)
	channel := make(chan string, bufferSize)

	go func() {
		defer close(channel)
		for msg := range r.pubsub.Channel() {
			if msg.Payload == key {
				channel <- msg.Payload
			}
		}
	}()

	return channel, nil
}

func (r *RedisPubSub) Publish(ctx context.Context, topic, message string) error {
	return r.client.Publish(ctx, topic, message).Err()
}

func (r *RedisPubSub) Close() error {
	if r.pubsub != nil {
		return r.pubsub.Close()
	}
	return nil
}

func (r *RedisPubSub) GetStrategy(_ context.Context) (string, error) {
	return r.strategy, nil
}
