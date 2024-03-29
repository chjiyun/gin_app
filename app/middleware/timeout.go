package middleware

import (
	"context"
	"errors"
	"gin_app/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetTimeout middleware wraps the request context with a timeout
func (m Middleware) SetTimeout() gin.HandlerFunc {
	return func(c *gin.Context) {

		timeout := time.Duration(config.Cfg.Server.Timeout) * time.Millisecond
		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			// check if context timeout was reached
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {

				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			//cancel to clear resources after finished
			cancel()
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
