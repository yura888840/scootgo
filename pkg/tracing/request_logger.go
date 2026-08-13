package tracing

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type HTTPRequest struct {
	Method        string        `json:"method"`
	RequestURI    string        `json:"request_uri"`
	StartTime     string        `json:"start_time"`
	LatencyInNano time.Duration `json:"latency_nano"`
	UserAgent     string        `json:"user_agent"`
	ClientIP      string        `json:"client_ip"`
	Status        int           `json:"status"`
	Size          int           `json:"size"`
}

func LogRequest(log *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		defer func() {
			req := c.Request
			logPayload := HTTPRequest{
				Method:        req.Method,
				RequestURI:    req.RequestURI,
				StartTime:     start.Format(time.RFC3339),
				LatencyInNano: time.Since(start),
				UserAgent:     req.UserAgent(),
				ClientIP:      c.ClientIP(),
				Status:        c.Writer.Status(),
				Size:          c.Writer.Size(),
			}

			log.Info().Ctx(c).Interface("req_data", logPayload).Msg("processed request")
		}()

		c.Next()
	}
}
