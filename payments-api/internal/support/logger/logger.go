package logger

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

func New(cfg config.Logger) (Logger, error) {
	switch cfg.Strategy {
	case "slog":
		return NewSlog(cfg)
	default:
		return nil, fmt.Errorf("router strategy not suported: %s", cfg.Strategy)
	}
}
