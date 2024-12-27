package router

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/docs"

	ginHandler "github.com/jtonynet/go-payments-api/internal/adapter/http/handler"
	ginMiddleware "github.com/jtonynet/go-payments-api/internal/adapter/http/middleware"
)

type Gin struct {
	app bootstrap.RESTApp
}

func NewGin(cfg config.Router, app bootstrap.RESTApp) (Router, error) {
	return Gin{app}, nil
}

func (gr Gin) HandleRequests(_ context.Context, cfg config.API) error {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	if cfg.MetricEnabled {
		initializeMetricsRoute(r, cfg)
	}

	v1 := r.Group("/")
	v1.Use(ginMiddleware.ConfigInject(cfg))
	v1.Use(ginMiddleware.AppInject(gr.app))

	v1.GET("/liveness", ginHandler.Liveness)
	v1.POST("/payment", ginHandler.PaymentExecution)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := fmt.Sprintf(":%s", cfg.Port)
	r.Run(port)

	return nil
}

func initializeMetricsRoute(r *gin.Engine, cfg config.API) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	ginMiddleware.InitPrometheus(r, cfg)
	r.Use(ginMiddleware.Prometheus(cfg))
}
