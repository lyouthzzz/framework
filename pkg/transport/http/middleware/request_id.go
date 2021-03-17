package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

const (
	requestIdKey = "X-Request-Id"
)

func RequestIdWith(c *gin.Context, requestId string) {
	c.Set(requestId, requestId)
}

func RequestMustFrom(c *gin.Context) string {
	val, _ := c.Get(requestIdKey)
	return val.(string)
}

func generateRequestId() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestId string
		if requestId = c.Request.Header.Get(requestIdKey); requestId == "" {
			requestId = generateRequestId()

		}
		RequestIdWith(c, requestId)
		c.Request.Header.Set(requestIdKey, requestId)

		c.Next()
	}
}
