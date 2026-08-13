package tracing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	testURI = "/"
)

var tests = []struct {
	name          string
	correlationID string
}{
	{
		"Without X-Correlation-Id, will be generated in the response",
		"",
	},
	{
		"With incorrect X-Correlation-Id, request will be aborted",
		"e32910d1-85sa6c-4b3f-ae10-sasa955141796aaaa",
	},
	{
		"With existing X-Correlation-Id",
		"e32910d1-856c-4b3f-ae10-2955141796ba",
	},
}

func TestSetCorrelationID(t *testing.T) {
	var err error
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(SetCorrelationID())
	router.GET(testURI, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})
	for _, test := range tests {
		responseRecorder := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(responseRecorder)

		testContext.Request, err = http.NewRequest("GET", testURI, http.NoBody)
		assert.NoError(t, err)

		if len(test.correlationID) > 0 {
			testContext.Request.Header.Set(ReqHeaderCorrelationID, test.correlationID)
		}
		router.HandleContext(testContext)
		value, isExists := testContext.Get(ReqHeaderCorrelationID)
		if isExists == true {
			if len(test.correlationID) > 0 {
				assert.Equal(t, test.correlationID, value)
			} else {
				assert.NotEmpty(t, value)
			}
		} else {
			assert.Equal(t, 400, responseRecorder.Code)
		}
	}
}
