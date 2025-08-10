package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// This Middleware adds a trace ID to the request context for tracking purposes.

func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = fmt.Sprintf("%X", time.Now().UnixNano())
			c.Request.Header.Set("X-Trace-ID", traceID)
		}

		c.Set("traceId", traceID)
		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Next()
	}
}
