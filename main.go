package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"micro-fiber-test/pkg/config"
	"micro-fiber-test/pkg/exceptions"
	endpoints "micro-fiber-test/pkg/handlers"
	"micro-fiber-test/pkg/middlewares"
	"micro-fiber-test/pkg/repository/impl"
	svcImpl "micro-fiber-test/pkg/service/impl"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const V1_ROOT = "/api/v1"
const ORG_V1_ROOT = V1_ROOT + "/organizations"
const ORG_V1_ORG_CODE = ORG_V1_ROOT + "/:orgCode"
const SECTORS_V1_ROOT = ORG_V1_ORG_CODE + "/sectors"
const SECTORS_V1_SECTOR_CODE = SECTORS_V1_ROOT + "/:sectorCode"
const USERS_V1_ROOT = ORG_V1_ORG_CODE + "/users"
const USERS_V1_USER_ID = USERS_V1_ROOT + "/:userId"

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
		AppName:           "micro-fiber-test",
		CaseSensitive:     true,
		StrictRouting:     true,
		EnablePrintRoutes: true,
		UnescapePath:      true,
		ErrorHandler:      defErrorHandler,
	}

	stdLogger.Info("Application -> Setup")
	app := fiber.New(fConfig)
	app.Use(middlewares.NewAccessLogger(accessLogger))
	app.Use(middlewares.NewHttpFilterLogger(stdLogger))
	if configuration.PrometheusEnabled {
		app.Use(configuration.PrometheusMetricsPath, basicauth.New(middlewares.NewBasicAuthConfig(configuration.BasicAuthUser, configuration.BasicAuthPass)))
	}
	if configuration.PrometheusEnabled {
		prometheus := fiberprometheus.New("micro-fiber-test")
		prometheus.RegisterAt(app, configuration.PrometheusMetricsPath)
		app.Use(prometheus.Middleware)
	}

	app.Static("/", "./static")

	// Organizations
	app.Get(ORG_V1_ROOT, endpoints.MakeOrgFindAll(configuration.TenantId, orgSvc))
	app.Post(ORG_V1_ROOT, endpoints.MakeOrgCreateEndpoint(configuration.RdbmsUrl, configuration.TenantId, orgSvc))
	app.Put(ORG_V1_ORG_CODE, endpoints.MakeOrgUpdateEndpoint(configuration.TenantId, orgSvc))
	app.Delete(ORG_V1_ORG_CODE, endpoints.MakeOrgDeleteEndpoint(configuration.TenantId, orgSvc))
	app.Get(ORG_V1_ORG_CODE, endpoints.MakeOrgFindByCodeEndpoint(configuration.TenantId, orgSvc))

	// Sectors
	app.Get(SECTORS_V1_ROOT, endpoints.MakeSectorsFindByOrga(configuration.TenantId, orgSvc, sectorSvc))
	app.Post(SECTORS_V1_ROOT, endpoints.MakeSectorCreateEndpoint(configuration.TenantId, orgSvc, sectorSvc))
	app.Put(SECTORS_V1_SECTOR_CODE, endpoints.MakeSectorUpdateEndpoint(configuration.TenantId, orgSvc, sectorSvc))
	app.Delete(SECTORS_V1_SECTOR_CODE, endpoints.MakeSectorDeleteEndpoint(configuration.TenantId, orgSvc, sectorSvc))

	// Users
	app.Get(USERS_V1_ROOT, endpoints.MakeUserSearchFilter(configuration.TenantId, userSvc, orgSvc))
	app.Get(USERS_V1_USER_ID, endpoints.MakeUserFindByCode(configuration.TenantId, userSvc, orgSvc))
	app.Post(USERS_V1_ROOT, endpoints.MakeUserCreateEndpoint(configuration.TenantId, userSvc, orgSvc))
	app.Put(USERS_V1_USER_ID, endpoints.MakeUserUpdate(configuration.TenantId, userSvc, orgSvc))
	app.Delete(USERS_V1_USER_ID, endpoints.MakeUserDelete(configuration.TenantId, userSvc, orgSvc))

	// OAuth and authentication
	app.Get("/api/v1/authenticate", endpoints.MakeGitlabAuthentication(store, configuration.OAuthGithub, configuration.OAuthClientId, configuration.OAuthRedirectUri))
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
