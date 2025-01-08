package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jtonynet/go-payments-api/config"
)

/*
	font: https://betterstack.com/community/guides/logging/logging-in-go/
*/

var levelNameToValue = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

type LokiHandler struct {
	client  *http.Client
	url     string
	options *slog.HandlerOptions
}

func NewLokiHandler(url string, opts *slog.HandlerOptions) *LokiHandler {
	return &LokiHandler{
		client:  &http.Client{Timeout: 10 * time.Second},
		url:     url,
		options: opts,
	}
}

func (h *LokiHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.options.Level.Level()
}

func (h *LokiHandler) Handle(ctx context.Context, r slog.Record) error {
	var buf bytes.Buffer

	stream := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": map[string]string{
					"level": r.Level.String(),
				},
				"values": [][]string{
					{
						fmt.Sprintf("%d", time.Now().UnixNano()),
						formatRecord(r),
					},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(stream); err != nil {
		return fmt.Errorf("failed to encode log entry: %w", err)
	}

	resp, err := h.client.Post(h.url, "application/json", &buf)
	if err != nil {
		return fmt.Errorf("failed to send log to Loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("loki responded with status: %s", resp.Status)
	}

	return nil
}

func (h *LokiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newOpts := *h.options
	newOpts.AddSource = h.options.AddSource
	return &LokiHandler{
		client:  h.client,
		url:     h.url,
		options: &newOpts,
	}
}

func (h *LokiHandler) WithGroup(name string) slog.Handler {
	return h
}

func formatRecord(r slog.Record) string {
	var attrs []string
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value))
		return true
	})
	return fmt.Sprintf("%s %s [%s]", r.Time.Format(time.RFC3339), r.Message, attrs)
}

type SLogger struct {
	instance *slog.Logger
}

func NewSlog(cfg config.Logger) (Logger, error) {
	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     levelNameToValue[cfg.Level],
	}

	var handler slog.Handler
	switch cfg.Output {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	case "loki":
		handler = NewLokiHandler(cfg.LokiPushURL, opts)
	default:
		return nil, fmt.Errorf("log strategy %s format: %s not suported", cfg.Strategy, cfg.Output)
	}

	instance := slog.New(handler)

	return &SLogger{
		instance: instance,
	}, nil
}

func (l SLogger) Info(msg string, args ...interface{}) {
	l.instance.Info(msg, args...)
}

func (l SLogger) Debug(msg string, args ...interface{}) {
	l.instance.Debug(msg, args...)
}

func (l SLogger) Warn(msg string, args ...interface{}) {
	l.instance.Warn(msg, args...)
}

func (l SLogger) Error(msg string, args ...interface{}) {
	l.instance.Error(msg, args...)
}
