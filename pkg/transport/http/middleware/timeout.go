package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func Timeout(t time.Duration) gin.HandlerFunc {
	if time.Duration(0) == 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return func(c *gin.Context) {
		timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer func() {
			if c.Err() == context.DeadlineExceeded {
				c.Abort()
			}
			cancel()
		}()
		c.Request = c.Request.WithContext(timeoutCtx)

		c.Next()
	}
}
