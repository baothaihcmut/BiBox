package middleware

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/monitor"
	"github.com/gin-gonic/gin"
)

func PrometheuseMiddleware(svc monitor.PrometheusService) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		method := c.Request.Method
		uri := c.FullPath()
		if uri == "" {
			uri = "unknown"
		}
		svc.RecordRequestDuration(c.Request.Context(), method, uri, status, duration)
		svc.IncRequestTotal(c.Request.Context(), method, uri, status)

	}
}
