package pubSub

import (
	"context"
	"fmt"
	"sync"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	client *redis.Client
	pubsub *redis.PubSub

	bufferSize            int
	subscriptionListeners sync.Map

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

	rps := &RedisPubSub{
		client:                client,
		strategy:              cfg.Strategy,
		bufferSize:            cfg.BufferSize,
		subscriptionListeners: sync.Map{},
	}

	_, err := rps.subscribe(context.Background())
	if err != nil {
		return &RedisPubSub{}, err
	}

	return rps, nil
}

func (r *RedisPubSub) Subscribe(_ context.Context, key string) (<-chan string, error) {
	if value, ok := r.subscriptionListeners.Load(key); ok {
		subscriptionListener, _ := value.(chan string)
		return subscriptionListener, nil
	}

	listnerBufferSize := 1
	listnerChannel := make(chan string, listnerBufferSize)
	r.subscriptionListeners.Store(key, listnerChannel)
	return listnerChannel, nil
}

func (r *RedisPubSub) UnSubscribe(_ context.Context, key string) error {
	r.subscriptionListeners.Delete(key)
	return nil
}

func (r *RedisPubSub) subscribe(_ context.Context) (<-chan string, error) {
	keyspaceChannel := fmt.Sprintf("__keyevent@%d__:expired", r.client.Options().DB)
	r.pubsub = r.client.Subscribe(context.Background(), keyspaceChannel)
	channel := make(chan string, r.bufferSize)

	go func() {
		for msg := range r.pubsub.Channel() {
			if listner, ok := r.subscriptionListeners.Load(msg.Payload); ok {
				listener, _ := listner.(chan string)
				listener <- msg.Payload
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
