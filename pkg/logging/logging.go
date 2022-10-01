package logging

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func New(cfg logger.Config, zapLogger *zap.Logger) fiber.Handler {

	var tzLoc *time.Location
	tz, err := time.LoadLocation(cfg.TimeZone)
	if err != nil || tz == nil {
		tzLoc = time.Local
	} else {
		tzLoc = tz
	}
	var timestamp atomic.Value
	timestamp.Store(time.Now().In(tzLoc).Format(cfg.TimeFormat))
	if strings.Contains(cfg.Format, "${time}") {
		go func() {
			for {
				time.Sleep(cfg.TimeInterval)
				timestamp.Store(time.Now().In(tzLoc).Format(cfg.TimeFormat))
			}
		}()
	}

	var (
		once sync.Once
		//mu         sync.Mutex
		errHandler fiber.ErrorHandler
	)

	return func(c *fiber.Ctx) (err error) {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Set error handler once
		once.Do(func() {
			// get longested possible path
			// override error handler
			errHandler = c.App().ErrorHandler
		})

		var start, stop time.Time

		start = time.Now()

		// Handle request, store err for logging
		chainErr := c.Next()

		stop = time.Now()
		duration := stop.Sub(start).Round(time.Millisecond)

		// Manually call error handler
		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		reqHeaders := make([]string, 0)
		for k, v := range c.GetReqHeaders() {
			reqHeaders = append(reqHeaders, k+"="+v)
		}
		zapLogger.Debug("HTTP",
			zap.Field{Key: "method", Type: zapcore.StringType, String: c.Method()},
			zap.Field{Key: "path", Type: zapcore.StringType, String: c.Path()},
			zap.Field{Key: "status", Type: zapcore.Int64Type, Integer: int64(c.Response().StatusCode())},
			zap.Field{Key: "ellapsed", Type: zapcore.Int64Type, Integer: duration.Milliseconds()},
			zap.Field{Key: "http.request.headers", Type: zapcore.StringType, String: strings.Join(reqHeaders, "&")},
			zap.Field{Key: "http.request.body", Type: zapcore.StringType, String: string(c.Body())},
			zap.Field{Key: "http.response.status", Type: zapcore.Int64Type, Integer: int64(c.Response().StatusCode())},
			zap.Field{Key: "http.response.body", Type: zapcore.StringType, String: string(c.Response().Body())},
		)
		return nil
	}

}
