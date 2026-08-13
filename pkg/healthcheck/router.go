package healthcheck

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Routes(
	router *gin.Engine,
	log *zerolog.Logger,
	isDBConnected func(context.Context) error,
	isAMQPConnected func() bool,
	errChan chan<- error,
) {
	router.GET("/monitoring/healthcheck/ping", healthcheckHandler)
	router.GET("/monitoring/ready", readyHandler(log, isDBConnected, isAMQPConnected, errChan))
}
