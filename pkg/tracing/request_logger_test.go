package tracing

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

const (
	testReqLoggerURI = "/"
)

type logHook struct {
	logEvents []zerolog.Event
}

func (logHook *logHook) Run(logEvent *zerolog.Event, _ zerolog.Level, _ string) {
	logHook.logEvents = append(logHook.logEvents, *logEvent)
}

func TestRequestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := zerolog.New(os.Stdout)
	logHook := &logHook{}
	logger = logger.Hook(logHook)

	router := gin.Default()
	router.Use(LogRequest(&logger))
	router.GET(testReqLoggerURI, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	req, err := http.NewRequest(http.MethodGet, testReqLoggerURI, http.NoBody)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	amountOfLogEvents := len(logHook.logEvents)
	assert.Equal(t, 1, amountOfLogEvents)
}
