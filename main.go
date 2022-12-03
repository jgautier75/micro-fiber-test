package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"micro-fiber-test/pkg/config"
	"micro-fiber-test/pkg/exceptions"
	"micro-fiber-test/pkg/handlers"
	"micro-fiber-test/pkg/middlewares"
	"micro-fiber-test/pkg/repository/impl"
	svcImpl "micro-fiber-test/pkg/service/impl"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {

	// Load config file
	configuration := config.LoadConfigFile("config/config.yaml")

	var kSql = koanf.New(".")
	errLoadSql := kSql.Load(file.Provider("config/sql_queries.toml"), toml.Parser())
	if errLoadSql != nil {
		panic(errLoadSql)
	}

	// Setup loggers
	accessLogger := configureLogger(configuration.LogsMetrics, false, false)
	stdLogger := configureLogger(configuration.LogsStd, true, true)

	dbPool, poolErr := configuration.SetupCnxPool(stdLogger)
	if poolErr != nil {
		panic(poolErr)
	}

	stdLogger.Info("Dao & Services -> Setup & inject")
	orgDao := impl.NewOrgDao(dbPool, kSql)
	sectorDao := impl.NewSectorDao(dbPool, kSql)
	userDao := impl.NewUserDao(dbPool, kSql)
	orgSvc := svcImpl.NewOrgService(dbPool, orgDao, sectorDao)
	sectorSvc := svcImpl.NewSectorService(sectorDao)
	userSvc := svcImpl.NewUserService(userDao)

	var defErrorHandler = func(c *fiber.Ctx, err error) error {
		var e *fiber.Error
		code := fiber.StatusInternalServerError
		if errors.As(err, &e) {
			code = e.Code
			if code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError {
				apiError := exceptions.ConvertToFunctionalError(err, code)
				return c.Status(code).JSON(apiError)
			} else {
				apiError := exceptions.ConvertToInternalError(err)
				return c.Status(code).JSON(apiError)
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(exceptions.ConvertToInternalError(err))
	}

	redisStorage := redis.New(redis.Config{
		Host:      configuration.RedisHost,
		Port:      configuration.RedisPort,
		Username:  configuration.RedisUser,
		Password:  configuration.RedisPass,
		URL:       "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	},
	)

	// Session storage in redis
	defCfg := session.ConfigDefault
	defCfg.Storage = redisStorage
	store := session.New(defCfg)
	fConfig := fiber.Config{
		CaseSensitive:     true,
		StrictRouting:     true,
		EnablePrintRoutes: false,
		UnescapePath:      true,
		ErrorHandler:      defErrorHandler,
	}

	stdLogger.Info("Application -> Setup")
	app := fiber.New(fConfig)
	app.Use(middlewares.NewAccessLogger(accessLogger))
	app.Use(middlewares.NewHttpFilterLogger(stdLogger))
	if configuration.PrometheusEnabled {
		prometheus := middlewares.PrometheusNew("micro-fiber-test")
		prometheus.RegisterAt(app, configuration.PrometheusMetricsPath)
		app.Use(prometheus.Middleware)
	}

	app.Static("/", "./static")

	// Organizations
	app.Get("/api/v1/organizations", endpoints.MakeOrgFindAll(configuration.TenantId, orgSvc))
	app.Post("/api/v1/organizations", endpoints.MakeOrgCreateEndpoint(configuration.RdbmsUrl, configuration.TenantId, orgSvc))
	app.Put("/api/v1/organizations/:orgCode", endpoints.MakeOrgUpdateEndpoint(configuration.TenantId, orgSvc))
	app.Delete("/api/v1/organizations/:orgCode", endpoints.MakeOrgDeleteEndpoint(configuration.TenantId, orgSvc))
	app.Get("/api/v1/organizations/:orgCode", endpoints.MakeOrgFindByCodeEndpoint(configuration.TenantId, orgSvc))

	// Sectors
	app.Get("/api/v1/organizations/:orgCode/sectors", endpoints.MakeSectorsFindByOrga(configuration.TenantId, orgSvc, sectorSvc))
	app.Post("/api/v1/organizations/:orgCode/sectors", endpoints.MakeSectorCreateEndpoint(configuration.TenantId, orgSvc, sectorSvc))
	app.Put("/api/v1/organizations/:orgCode/sectors/:sectorCode", endpoints.MakeSectorUpdateEndpoint(configuration.TenantId, orgSvc, sectorSvc))
	app.Delete("/api/v1/organizations/:orgCode/sectors/:sectorCode", endpoints.MakeSectorDeleteEndpoint(configuration.TenantId, orgSvc, sectorSvc))

	// Users
	app.Get("/api/v1/organizations/:orgCode/users", endpoints.MakeUserSearchFilter(configuration.TenantId, userSvc, orgSvc))
	app.Get("/api/v1/organizations/:orgCode/users/:userId", endpoints.MakeUserFindByCode(configuration.TenantId, userSvc, orgSvc))
	app.Post("/api/v1/organizations/:orgCode/users", endpoints.MakeUserCreateEndpoint(configuration.TenantId, userSvc, orgSvc))
	app.Put("/api/v1/organizations/:orgCode/users/:userId", endpoints.MakeUserUpdate(configuration.TenantId, userSvc, orgSvc))
	app.Delete("/api/v1/organizations/:orgCode/users/:userId", endpoints.MakeUserDelete(configuration.TenantId, userSvc, orgSvc))

	// OAuth and authentication
	app.Get("/api/v1/authenticate", endpoints.MakeGitlabAuthentication(store, configuration.OAuthGitlab, configuration.OAuthClientId, configuration.OAuthRedirectUri))
	app.Get("/oauth/redirect", endpoints.MakeOAuthAuthorize(store, configuration.OAuthCallbackUrl, configuration.OAuthClientId, configuration.OAuthClientSecret, configuration.OAuthDebug))

	go func() {
		stdLogger.Info("Application -> ListenTLS")
		if errTls := app.ListenTLS(":"+configuration.ServerPort, "cert.pem", "key.pem"); errTls != nil {
			panic(errTls)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	stdLogger.Info("Closing connections pool")
	dbPool.Close()

}

// Configure metrics logger
func configureLogger(logCfg string, consoleOutput bool, callerAndStack bool) *zap.Logger {
	zapConfig := zap.NewProductionEncoderConfig()
	zapConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
	}
	fileEncoder := zapcore.NewJSONEncoder(zapConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(zapConfig)
	logFile, _ := os.OpenFile(logCfg, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.InfoLevel
	var core zapcore.Core
	if consoleOutput {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel))
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		)
	}
	if callerAndStack {
		return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		return zap.New(core)
	}
}
