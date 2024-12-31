package pubSub

import (
	"context"
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
)

type PubSub interface {
	GetStrategy(ctx context.Context) (string, error)
	Subscribe(ctx context.Context) (<-chan string, error)
	GetSubscription(ctx context.Context) (<-chan string, error)
	Publish(ctx context.Context, topic, message string) error
	Close() error
}

func New(cfg config.PubSub) (PubSub, error) {
	switch cfg.Strategy {
	case "redis":
		return NewRedisPubSub(cfg)
	default:
		return nil, fmt.Errorf("pubsub strategy not suported: %s", cfg.Strategy)
	}
}
