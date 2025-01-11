package logger

import (
	"context"
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
)

type contextKey string

const (
	CtxResponseCodeKey   contextKey = "code"
	CtxExecutionTimeKey  contextKey = "execution_time_in_ms"
	CtxTransactionUIDKey contextKey = "transaction_uid"
	CtxAccountUIDKey     contextKey = "account_uid"
)

var CtxKeysMap = map[string]contextKey{
	"code":                 CtxResponseCodeKey,
	"execution_time_in_ms": CtxExecutionTimeKey,
	"transaction_uid":      CtxTransactionUIDKey,
	"account_uid":          CtxAccountUIDKey,
}

type Logger interface {
	Info(ctx context.Context, msg string, args ...interface{})
	Debug(ctx context.Context, msg string, args ...interface{})
	Warn(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, args ...interface{})
}

func New(cfg config.Logger) (Logger, error) {
	switch cfg.Strategy {
	case "slog":
		return NewSlog(cfg)
	default:
		return nil, fmt.Errorf("router strategy not suported: %s", cfg.Strategy)
	}
}
