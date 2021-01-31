package http

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/lyouthzzz/framework/transport/http/middleware"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
)

type Option struct {
	Port int
	Mode string
}

type InitRouters func(r *gin.Engine)

func NewOption(port int, mode string) *Option {
	return &Option{Port: port, Mode: mode}
}

func NewRouter(opt *Option, logger *zerolog.Logger, initRouters InitRouters, tracer opentracing.Tracer) *gin.Engine {
	if opt.Mode == "" {
		opt.Mode = gin.DebugMode
	}

	gin.SetMode(opt.Mode)

	r := gin.New()
	ginpprof.Wrap(r)

	r.Use(middleware.RequestId())
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger))
	r.Use(ginhttp.Middleware(tracer))

	initRouters(r)

	return r
}
