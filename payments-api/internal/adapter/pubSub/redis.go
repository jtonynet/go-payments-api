package pubSub

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	client *redis.Client
	pubsub *redis.PubSub

	subscriptionListeners sync.Map
	strategy              string
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
		subscriptionListeners: sync.Map{},
	}

	err := rps.subscribe(context.Background())
	if err != nil {
		return &RedisPubSub{}, err
	}

	return rps, nil
}

func (r *RedisPubSub) Subscribe(_ context.Context, key string) (<-chan string, error) {
	if listner, ok := r.subscriptionListeners.Load(key); ok {
		subscriptionListener, _ := listner.(chan string)
		return subscriptionListener, nil
	}

	listnerBufferSize := 1
	listnerChannel := make(chan string, listnerBufferSize)
	r.subscriptionListeners.Store(key, listnerChannel)
	return listnerChannel, nil
}

func (r *RedisPubSub) UnSubscribe(_ context.Context, key string) error {
	if listener, ok := r.subscriptionListeners.Load(key); ok {
		r.subscriptionListeners.Delete(key)
		close(listener.(chan string))
	}

	return nil
}

func (r *RedisPubSub) subscribe(_ context.Context) error {
	keyspaceChannel := fmt.Sprintf("__keyevent@%d__:expired", r.client.Options().DB)
	r.pubsub = r.client.Subscribe(context.Background(), keyspaceChannel)

	go func() {
		for msg := range r.pubsub.Channel() {
			// Necessary to message multiple `transaction` requests for the same `account` by `UIDs`
			keyPrefix := msg.Payload // msg.Payload is a key `accountUID` of expired register
			r.subscriptionListeners.Range(func(key, value any) bool {
				if strKey, ok := key.(string); ok && strings.HasPrefix(strKey, keyPrefix) {
					if ch, ok := value.(chan string); ok {
						ch <- keyPrefix
					}
				}
				return true
			})
		}
	}()

	return nil
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
