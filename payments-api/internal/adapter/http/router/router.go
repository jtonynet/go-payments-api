package router

import (
	"context"
	"fmt"

	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/config"
)

type Router interface {
	HandleRequests(ctx context.Context, cfg config.API) error
}

func New(cfg config.Router, app bootstrap.RESTApp) (Router, error) {
	switch cfg.Strategy {
	case "gin":
		return NewGin(cfg, app)
	default:
		return nil, fmt.Errorf("router strategy not suported: %s", cfg.Strategy)
	}
}
