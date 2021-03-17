package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Recovery(logger *zerolog.Logger) gin.HandlerFunc {
	return gin.RecoveryWithWriter(logger)
}
