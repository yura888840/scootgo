package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yura888840/scootgo/pkg/config"
)

func TestRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	configs := map[string]string{config.ConfigAppName: "boilerplate-service"}

	NewRouter().Routes(router, configs)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)

	var body map[string]string
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &body))
	assert.Equal(t, "ok", body["status"])
	assert.Equal(t, "boilerplate-service", body["service"])
}
