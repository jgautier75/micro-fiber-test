package logging

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"sync"
	"time"
)

func New(zapLogger *zap.Logger) fiber.Handler {

	var (
		once       sync.Once
		errHandler fiber.ErrorHandler
	)

	return func(c *fiber.Ctx) (err error) {

		defer func(zapLogger *zap.Logger) {
			_ = zapLogger.Sync()
		}(zapLogger)

		// Set error handler once
		once.Do(func() {
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

		var logReq = !strings.Contains(c.Path(), "assets") && !strings.Contains(c.Path(), ".html")
		if logReq {
			zapLogger.Debug("HTTP",
				zap.Field{Key: "method", Type: zapcore.StringType, String: c.Method()},
				zap.Field{Key: "path", Type: zapcore.StringType, String: c.Path()},
				zap.Field{Key: "ellapsed", Type: zapcore.Int64Type, Integer: duration.Milliseconds()},
				zap.Field{Key: "http.response.status", Type: zapcore.Int64Type, Integer: int64(c.Response().StatusCode())},
			)
		}
		return nil
	}

}
