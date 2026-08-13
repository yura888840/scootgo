package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var errTestChannelFull = errors.New("channel is full")

func testLogger() *zerolog.Logger {
	log := zerolog.Nop()
	return &log
}

func setup(isDBConnected func(context.Context) error, isAMQPConnected func() bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Routes(router, testLogger(), isDBConnected, isAMQPConnected, nil)

	return router
}

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name          string
		uri           string
		responseCode  int
		response      string
		dbErr         error
		amqpConnected bool
	}{
		{
			"health check",
			"/monitoring/healthcheck/ping",
			200,
			"{\"status\":\"ok\"}",
			nil,
			true,
		},
		{
			"readiness check OK",
			"/monitoring/ready",
			200,
			"{\"dbReady\":true,\"amqpReady\":true}",
			nil,
			true,
		},
		{
			"readiness check DB fail",
			"/monitoring/ready",
			503,
			"{\"dbReady\":false,\"amqpReady\":true}",
			assert.AnError,
			true,
		},
		{
			"readiness check Rabbit fail",
			"/monitoring/ready",
			503,
			"{\"dbReady\":true,\"amqpReady\":false}",
			nil,
			false,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consecutiveReadinessFailures.Store(0)

			isDBConnected := func(_ context.Context) error {
				return tt.dbErr
			}
			isAMQPConnected := func() bool {
				return tt.amqpConnected
			}
			router := setup(isDBConnected, isAMQPConnected)

			req, err := http.NewRequestWithContext(ctx, "GET", tt.uri, http.NoBody)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.responseCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.response)
		})
	}
}

func TestReadinessFailureThresholdTriggersRestartError(t *testing.T) {
	consecutiveReadinessFailures.Store(0)
	errChan := make(chan error, 1)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Routes(router, testLogger(), func(_ context.Context) error { return assert.AnError }, func() bool { return true }, errChan)

	ctx := context.Background()
	for i := 0; i < 2; i++ {
		req, err := http.NewRequestWithContext(ctx, "GET", "/monitoring/ready", http.NoBody)
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Len(t, errChan, 0)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "/monitoring/ready", http.NoBody)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Len(t, errChan, 1)
}

func TestReadinessFailureCounterResetsOnSuccess(t *testing.T) {
	consecutiveReadinessFailures.Store(0)
	errChan := make(chan error, 1)

	var dbErr error
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Routes(router, testLogger(), func(_ context.Context) error { return dbErr }, func() bool { return true }, errChan)

	ctx := context.Background()

	dbErr = assert.AnError
	for i := 0; i < 2; i++ {
		req, err := http.NewRequestWithContext(ctx, "GET", "/monitoring/ready", http.NoBody)
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	}

	dbErr = nil
	req, err := http.NewRequestWithContext(ctx, "GET", "/monitoring/ready", http.NoBody)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, errChan, 0)

	dbErr = assert.AnError
	for i := 0; i < 2; i++ {
		req, reqErr := http.NewRequestWithContext(ctx, "GET", "/monitoring/ready", http.NoBody)
		assert.NoError(t, reqErr)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Len(t, errChan, 0)
	}

	req, err = http.NewRequestWithContext(ctx, "GET", "/monitoring/ready", http.NoBody)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Len(t, errChan, 1)
}

func TestHandleReadinessFailure_AtThresholdWithNilChannel(t *testing.T) {
	consecutiveReadinessFailures.Store(readinessFailureThreshold - 1)

	handleReadinessFailure(testLogger(), nil)

	assert.Equal(t, readinessFailureThreshold, consecutiveReadinessFailures.Load())
}

func TestHandleReadinessFailure_FullChannel(t *testing.T) {
	consecutiveReadinessFailures.Store(0)
	errChan := make(chan error, 1)
	errChan <- errTestChannelFull

	for i := int32(0); i < readinessFailureThreshold; i++ {
		handleReadinessFailure(testLogger(), errChan)
	}

	assert.Equal(t, readinessFailureThreshold, consecutiveReadinessFailures.Load())
	assert.Len(t, errChan, 1)
}
