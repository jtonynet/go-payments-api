package ginStrategy

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/docs"

	ginHandler "github.com/jtonynet/go-payments-api/internal/adapter/http/handler"
	ginMiddleware "github.com/jtonynet/go-payments-api/internal/adapter/http/middleware"

	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type Gin struct {
	app bootstrap.App
}

func New(cfg config.Router, app bootstrap.App) (port.Router, error) {
	return Gin{app}, nil
}

func (gr Gin) HandleRequests(_ context.Context, cfg config.API) error {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	apiGroup := r.Group("/")
	apiGroup.Use(ginMiddleware.ConfigInject(cfg))
	apiGroup.Use(ginMiddleware.AppInject(gr.app))

	apiGroup.GET("/liveness", ginHandler.Liveness)
	apiGroup.POST("/payment", ginHandler.PaymentExecution)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := fmt.Sprintf(":%s", cfg.Port)
	r.Run(port)

	return nil
}
