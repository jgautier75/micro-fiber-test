package middlewares

import (
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const BsTxId = "bsTxId"
const DebugHeader = "x-log-debug"

func NewHttpFilterLogger(zapLogger *zap.Logger) fiber.Handler {

	var (
		once       sync.Once
		errHandler fiber.ErrorHandler
	)

	return func(c *fiber.Ctx) (err error) {

		var bsTxId string
		bsTxId = uuid.New().String()
		c.Locals(BsTxId, bsTxId)
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

		// Handle request, store err for middlewares
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
		debugReq := false
		reqHeaders := make([]string, 0)
		for k, v := range c.GetReqHeaders() {
			if strings.ToLower(k) == DebugHeader && v[0] == "1" {
				debugReq = true
			}
			reqHeaders = append(reqHeaders, k+"="+v[0])
		}

		if logReq && debugReq {
			zapLogger.Info("HTTP",
				zap.Field{Key: "method", Type: zapcore.StringType, String: c.Method()},
				zap.Field{Key: "path", Type: zapcore.StringType, String: c.Path()},
				zap.Field{Key: "elapsed", Type: zapcore.Int64Type, Integer: duration.Milliseconds()},
				zap.Field{Key: "http.response.status", Type: zapcore.Int64Type, Integer: int64(c.Response().StatusCode())},
				zap.Field{Key: "http.request.headers", Type: zapcore.StringType, String: strings.Join(reqHeaders, "&")},
				zap.Field{Key: "http.request.body", Type: zapcore.StringType, String: string(c.Body())},
				zap.Field{Key: "http.response.body", Type: zapcore.StringType, String: string(c.Response().Body())},
				zap.Field{Key: "bsTransactionId", Type: zapcore.StringType, String: bsTxId},
			)
		}
		return nil
	}

}
