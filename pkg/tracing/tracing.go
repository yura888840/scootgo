package tracing

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const ReqHeaderCorrelationID = "X-Correlation-Id"

type TracingHook struct{}

func (h TracingHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	ctx := e.GetCtx()
	correlationID, ok := ctx.Value(ReqHeaderCorrelationID).(string)
	if ok {
		e.Str(strings.ToLower(ReqHeaderCorrelationID), correlationID)
	}
}

func SetCorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := ""
		if c.Request.Header.Get(ReqHeaderCorrelationID) == "" {
			correlationID = uuid.New().String()
		} else {
			reqUUID, err := uuid.Parse(c.Request.Header.Get(ReqHeaderCorrelationID))
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			} else {
				correlationID = reqUUID.String()
			}
		}
		c.Set(ReqHeaderCorrelationID, correlationID)
		c.Writer.Header().Add(ReqHeaderCorrelationID, correlationID)
		c.Writer.Header().Add("Access-Control-Expose-Headers", ReqHeaderCorrelationID)
		c.Next()
	}
}
