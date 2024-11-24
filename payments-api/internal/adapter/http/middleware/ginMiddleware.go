package ginMiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/config"
)

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
