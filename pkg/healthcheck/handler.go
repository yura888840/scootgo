package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const readinessFailureThreshold int32 = 3

var consecutiveReadinessFailures atomic.Int32
var errReadinessProbeFailedConsecutively = errors.New("readiness probe failed 3 consecutive times")

func handleReadinessFailure(log *zerolog.Logger, errChan chan<- error) {
	failureCount := consecutiveReadinessFailures.Add(1)
	if failureCount < readinessFailureThreshold || errChan == nil {
		return
	}

	select {
	case errChan <- errReadinessProbeFailedConsecutively:
		log.Error().Int32("failureCount", failureCount).Msg("readiness failure threshold reached, triggering restart")
	default:
		log.Warn().Int32("failureCount", failureCount).Msg("readiness failure threshold reached but error channel is full")
	}
}

func healthcheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type readinessResponse struct {
	DBReady   bool `json:"dbReady"`
	AMQPReady bool `json:"amqpReady"`
}

func readyHandler(
	log *zerolog.Logger,
	isDBConnected func(context.Context) error,
	isAMQPConnected func() bool,
	errChan chan<- error,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		response := readinessResponse{
			DBReady:   true,
			AMQPReady: true,
		}

		if err := isDBConnected(ctx); err != nil {
			response.DBReady = false
			log.Error().Err(err).Msg("database connection check failed")
		}

		if isAMQPConnected == nil || !isAMQPConnected() {
			response.AMQPReady = false
			log.Error().Msg("AMQP connection check failed")
		}

		statusCode := http.StatusOK
		if !response.DBReady || !response.AMQPReady {
			statusCode = http.StatusServiceUnavailable
			handleReadinessFailure(log, errChan)
		} else {
			consecutiveReadinessFailures.Store(0)
		}

		c.JSON(statusCode, response)
	}
}
