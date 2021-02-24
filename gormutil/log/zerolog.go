package log

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	glog "gorm.io/gorm/logger"
	"time"
)

func SlowThresholdOption(duration time.Duration) Option {
	return func(l *Logger) {
		l.slowThreshold = duration
	}
}

type Option func(l *Logger)

type Logger struct {
	slowThreshold time.Duration
}

func New(opts ...Option) *Logger {
	logger := &Logger{
		slowThreshold: 100 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

func (l Logger) LogMode(level glog.LogLevel) glog.Interface {
	return l
}

func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	zerolog.Ctx(ctx).Info().Msg(fmt.Sprintf(msg, args...))
}

func (l Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	zerolog.Ctx(ctx).Warn().Msg(fmt.Sprintf(msg, args...))
}

func (l Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	zerolog.Ctx(ctx).Warn().Msg(fmt.Sprintf(msg, args...))
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	zLog := zerolog.Ctx(ctx)
	var event *zerolog.Event

	if err != nil {
		event = zLog.Debug()
	} else {
		event = zLog.Trace()
	}

	var durKey string

	switch zerolog.DurationFieldUnit {
	case time.Nanosecond:
		durKey = "elapsed_ns"
	case time.Microsecond:
		durKey = "elapsed_us"
	case time.Millisecond:
		durKey = "elapsed_ms"
	case time.Second:
		durKey = "elapsed"
	case time.Minute:
		durKey = "elapsed_min"
	case time.Hour:
		durKey = "elapsed_hr"
	default:
		durKey = "elapsed_"
	}

	duration := time.Since(begin)

	event.Dur(durKey, duration)

	sql, rows := fc()
	if sql != "" {
		event.Str("sql", sql)
	}
	if rows > -1 {
		event.Int64("rows", rows)
	}
	if duration > l.slowThreshold {
		event.Bool("slow", true)
	}

	event.Send()
	return
}
