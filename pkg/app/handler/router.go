package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yura888840/scootgo/pkg/config"
)

type Router struct{}

func NewRouter() Router {
	return Router{}
}

func (r Router) Routes(engine *gin.Engine, configs map[string]string) {
	engine.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": configs[config.ConfigAppName],
		})
	})
}
