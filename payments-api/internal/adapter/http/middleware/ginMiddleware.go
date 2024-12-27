package ginMiddleware

import (
	"bytes"
	"encoding/json"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"

	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

var (
	startTime       = time.Now()
	transactionCode = "Not Applicable"

	// TIMER
	processUptime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_uptime_seconds",
			Help: "Total uptime of the process in seconds",
		},
		[]string{},
	)

	// CPU
	systemCPUUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_usage",
			Help: "System CPU usage percent",
		},
		[]string{"host"},
	)

	processCPUUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_cpu_usage",
			Help: "Process CPU usage percent",
		},
		[]string{"host"},
	)

	systemLoadAverage1m = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_load_average_1m",
			Help: "System load average over the last 1 minute",
		},
		[]string{"host"},
	)

	systemCPUCores = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_cores",
			Help: "Number of CPU cores on the system",
		},
		[]string{"host"},
	)

	// MEMORY
	memoryAlloc = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_alloc_bytes",
			Help: "Memory alloc of the process in seconds",
		},
		[]string{"host"},
	)

	memoryTotalAlloc = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_total_alloc_bytes",
			Help: "Memory total alloc of the process in seconds",
		},
		[]string{"host"},
	)

	processMemoryUsed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_used_bytes",
			Help: "Memory used by the process in bytes",
		},
		[]string{"host"},
	)

	processMemoryLimit = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_limit_bytes",
			Help: "Memory limit of the process in bytes",
		},
		[]string{"host"},
	)

	// HTTP REQUESTS
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gin",
			Subsystem: "gonic",
			Name:      "requests_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"statusCode", "transactionCode", "method", "host", "path"},
	)

	requestsDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "gin",
		Subsystem: "gonic",
		Name:      "request_duration",
		Help:      "The HTTP request latency bucket.",
	}, []string{"statusCode", "transactionCode", "method", "host", "path"})

	maxDurations        = make(map[string]float64)
	requestsMaxDuration = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "gin",
			Subsystem: "gonic",
			Name:      "requests_seconds_max",
			Help:      "The maximum HTTP request latency.",
		},
		[]string{"statusCode", "transactionCode", "method", "host", "path"},
	)
)

func InitPrometheus(r *gin.Engine, cfg config.API) {
	// Use OpenTelemetry middleware
	r.Use(otelgin.Middleware(cfg.Name))

	prometheus.Register(processUptime)

	prometheus.Register(systemCPUUsage)
	prometheus.Register(processCPUUsage)
	prometheus.Register(systemLoadAverage1m)
	prometheus.Register(systemCPUCores)

	prometheus.Register(memoryAlloc)
	prometheus.Register(memoryTotalAlloc)

	prometheus.Register(requestsTotal)
	prometheus.Register(requestsDuration)
	prometheus.Register(requestsMaxDuration)

	metricsLoop(cfg)
}

func metricsLoop(cfg config.API) {
	go func() {
		for {
			uptime := time.Since(startTime).Seconds()
			processUptime.WithLabelValues().Set(uptime)

			systemCPUPercent, err := cpu.Percent(0, true)
			if err == nil {
				systemCPUUsage.WithLabelValues(cfg.RestHost).Set(systemCPUPercent[0])
			}

			processCPUPercent, err := cpu.Percent(0, false)
			if err == nil {
				processCPUUsage.WithLabelValues(cfg.RestHost).Set(processCPUPercent[0])
			}

			loadAvg, err := load.Avg()
			if err == nil {
				systemLoadAverage1m.WithLabelValues(cfg.RestHost).Set(loadAvg.Load1)
			}

			cpuInfo, err := cpu.Info()
			if err == nil {
				systemCPUCores.WithLabelValues(cfg.RestHost).Set(float64(len(cpuInfo)))
			}

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			memoryAlloc.WithLabelValues(cfg.RestHost).Set(float64(m.Alloc))
			memoryTotalAlloc.WithLabelValues(cfg.RestHost).Set(float64(m.TotalAlloc))

			memInfo, err := mem.VirtualMemory()
			if err == nil {
				processMemoryUsed.WithLabelValues(cfg.RestHost).Set(float64(memInfo.Used))
				processMemoryLimit.WithLabelValues(cfg.RestHost).Set(float64(memInfo.Total))
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Prometheus(cfg config.API) gin.HandlerFunc {
	return func(c *gin.Context) {
		initialRequestTime := time.Now()

		writer := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		c.Next()

		if c.FullPath() == cfg.TransactionPath {
			bodyBytes := writer.body.Bytes()
			var responseBody map[string]interface{}
			if len(bodyBytes) > 0 {
				if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
					transactionCode = port.CODE_REJECTED_GENERIC
				}
			}

			if val, exists := responseBody["code"]; exists {
				if strVal, ok := val.(string); ok {
					transactionCode = strVal
				}
			}
		}

		statusCode := strconv.Itoa(c.Writer.Status())
		labels := []string{statusCode, transactionCode, c.Request.Method, c.Request.Host, c.FullPath()}

		requestsTotal.WithLabelValues(labels...).Inc()

		duration := time.Since(initialRequestTime).Seconds()
		requestsDuration.WithLabelValues(labels...).Observe(duration)

		maxDuration, exists := maxDurations[strings.Join(labels, "_")]
		if !exists || duration > maxDuration {
			maxDurations[strings.Join(labels, "_")] = duration
			requestsMaxDuration.WithLabelValues(labels...).Set(duration)
		}
	}
}

func AppInject(app bootstrap.RESTApp) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("app", app)
		c.Next()
	}
}

func ConfigInject(cfg config.API) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cfg", cfg)
		c.Next()
	}
}
