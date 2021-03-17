package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http/httputil"
)

func RequestDump(verbose bool, logger *zerolog.Logger) gin.HandlerFunc {
	if verbose {
		return func(c *gin.Context) {
			data, _ := httputil.DumpRequest(c.Request, false)
			logger.Info().Msgf("REQUEST DUMP\n: %s", data)

			c.Next()
		}
	}

	return func(c *gin.Context) {
		c.Next()
	}
}
