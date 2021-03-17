package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"strconv"
	"sync"
	"time"
)

const loggerKey = "zero_logger"

func LoggerWith(c *gin.Context, logger *zerolog.Logger) {
	c.Set(loggerKey, logger)
}

func LoggerMustFrom(c *gin.Context) *zerolog.Logger {
	val, _ := c.Get(loggerKey)
	return val.(*zerolog.Logger)
}

func Logger(logger *zerolog.Logger) gin.HandlerFunc {
	pool := &sync.Pool{
		New: func() interface{} {
			buf := new(bytes.Buffer)
			return buf
		},
	}

	return func(c *gin.Context) {
		start := time.Now()
		requestId := c.Request.Header.Get(requestIdKey)
		loggerWithId := logger.With().Str(requestId, requestId).Logger()

		c.Request = c.Request.WithContext(loggerWithId.WithContext(c.Request.Context()))

		LoggerWith(c, &loggerWithId)

		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path += "?" + c.Request.URL.RawQuery
		}

		c.Next()

		// [GIN] | 2021/01/18 - 23:59:33 | 200 | 2.28933742s | 183.194.242.255 | POST | /public/api/bigdata/ai/11?appkey=1c1cb632b9404dd7b6160934c7fe43c2 | b3d5c293e632c9612d64e7d11d25fcc4
		latency := time.Now().Sub(start).Seconds()
		w := pool.Get().(*bytes.Buffer)
		w.Reset()
		w.WriteString("[GIN]")
		w.WriteString(" | ")
		w.WriteString(time.Now().Format("2006/01/02 - 15:04:05"))
		w.WriteString(" | ")
		w.WriteString(strconv.Itoa(c.Writer.Status()))
		w.WriteString(" | ")
		w.WriteString(fmt.Sprintf("%f", latency))
		w.WriteString(" | ")
		w.WriteString(c.ClientIP())
		w.WriteString(" | ")
		w.WriteString(c.Request.Method)
		w.WriteString(" | ")
		w.WriteString(path)

		if c.Errors != nil && len(c.Errors) != 0 {
			w.WriteString(" | ")
			w.WriteString(c.Errors.String())
		}
		w.WriteString("\n")

		_, _ = w.WriteTo(loggerWithId)
		pool.Put(w)
	}
}
