package middlewares

import (
	"app/base/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// setup logging middleware
// ensures logging line after each http response with fields:
// duration_ms, status, userAgent, method, remoteAddr, url, param_*
func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		tStart := time.Now()
		c.Next()
		var fields []interface{}

		duration := time.Since(tStart).Nanoseconds() / 1e6
		fields = append(fields, "durationMs", duration,
			"status_code", c.Writer.Status(),
			"user_agent", c.Request.UserAgent(),
			"method", c.Request.Method,
			"remote_addr", c.Request.RemoteAddr,
			"url", c.Request.URL.String(),
			"content_encoding", c.Writer.Header().Get("Content-Encoding"),
			"account", c.GetInt(KeyAccount))

		for _, param := range c.Params {
			fields = append(fields, "param_"+param.Key, param.Value)
		}
		fields = append(fields, "request")

		if c.Writer.Status() < http.StatusInternalServerError {
			utils.LogInfo(fields...)
		} else {
			utils.LogError(fields...)
		}

		utils.ObserveSecondsSince(tStart, requestDurations.
			WithLabelValues(c.Request.Method+c.FullPath()))
	}
}
