package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"micro-fiber-test/pkg/exceptions"
	"micro-fiber-test/pkg/handlers"
	"micro-fiber-test/pkg/middlewares"
	"micro-fiber-test/pkg/repository/impl"
	svcImpl "micro-fiber-test/pkg/service/impl"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func main() {

	// Load config file
	var kConfig = koanf.New(".")
	errLoadCfg := kConfig.Load(file.Provider("config/config.yaml"), yaml.Parser())
	if errLoadCfg != nil {
		panic(errLoadCfg)
	}
	targetPort := kConfig.String("http.server.port")
	defaultTenantId := kConfig.Int64("app.tenant")
	accessLogFile := kConfig.String("app.accessLogFile")
	stdLogFile := kConfig.String("app.stdLogFile")
	clientId := kConfig.String("app.oauthClientId")
	clientSecret := kConfig.String("app.oauthClientSecret")
	oauthCallback := kConfig.String("app.oauthCallback")
	oauthRedirectUri := kConfig.String("app.oauthRedirectUri")
	oauthGitlab := kConfig.String("app.oauthGitlab")
	oauthDebug := kConfig.Bool("app.oauthDebug")
	dbUrl := kConfig.String("app.pgUrl")
	redisHost := kConfig.String("app.redisHost")
	redisStrPort := kConfig.String("app.redisPort")
	redisPort, errRedis := strconv.Atoi(redisStrPort)
	redisUser := kConfig.String("app.redisUser")
	redisPass := kConfig.String("app.redisPass")
	metricsPath := kConfig.String("app.metricsPath")
	prometheusEnabled := kConfig.Bool("app.prometheusEnabled")

	var kSql = koanf.New(".")
	errLoadSql := kSql.Load(file.Provider("config/sql_queries.toml"), toml.Parser())
	if errLoadSql != nil {
		panic(errLoadSql)
	}

	fmt.Printf("Redis port error [%v]", errRedis)

	// Setup loggers
	accessLogger := configureLogger(accessLogFile, false, false)
	stdLogger := configureLogger(stdLogFile, true, true)

	dbPool, poolErr := configureCnxPool(kConfig, stdLogger)
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
		Host:      redisHost,
		Port:      redisPort,
		Username:  redisUser,
		Password:  redisPass,
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
	if prometheusEnabled {
		prometheus := middlewares.PrometheusNew("micro-fiber-test")
		prometheus.RegisterAt(app, metricsPath)
		app.Use(prometheus.Middleware)
	}

	app.Static("/", "./static")

	// Organizations
	app.Get("/api/v1/organizations", endpoints.MakeOrgFindAll(defaultTenantId, orgSvc))
	app.Post("/api/v1/organizations", endpoints.MakeOrgCreateEndpoint(dbUrl, defaultTenantId, orgSvc))
	app.Put("/api/v1/organizations/:orgCode", endpoints.MakeOrgUpdateEndpoint(defaultTenantId, orgSvc))
	app.Delete("/api/v1/organizations/:orgCode", endpoints.MakeOrgDeleteEndpoint(defaultTenantId, orgSvc))
	app.Get("/api/v1/organizations/:orgCode", endpoints.MakeOrgFindByCodeEndpoint(defaultTenantId, orgSvc))

	// Sectors
	app.Get("/api/v1/organizations/:orgCode/sectors", endpoints.MakeSectorsFindByOrga(defaultTenantId, orgSvc, sectorSvc))
	app.Post("/api/v1/organizations/:orgCode/sectors", endpoints.MakeSectorCreateEndpoint(defaultTenantId, orgSvc, sectorSvc))
	app.Put("/api/v1/organizations/:orgCode/sectors/:sectorCode", endpoints.MakeSectorUpdateEndpoint(defaultTenantId, orgSvc, sectorSvc))
	app.Delete("/api/v1/organizations/:orgCode/sectors/:sectorCode", endpoints.MakeSectorDeleteEndpoint(defaultTenantId, orgSvc, sectorSvc))

	// Users
	app.Get("/api/v1/organizations/:orgCode/users", endpoints.MakeUserSearchFilter(defaultTenantId, userSvc, orgSvc))
	app.Get("/api/v1/organizations/:orgCode/users/:userId", endpoints.MakeUserFindByCode(defaultTenantId, userSvc, orgSvc))
	app.Post("/api/v1/organizations/:orgCode/users", endpoints.MakeUserCreateEndpoint(defaultTenantId, userSvc, orgSvc))
	app.Put("/api/v1/organizations/:orgCode/users/:userId", endpoints.MakeUserUpdate(defaultTenantId, userSvc, orgSvc))
	app.Delete("/api/v1/organizations/:orgCode/users/:userId", endpoints.MakeUserDelete(defaultTenantId, userSvc, orgSvc))

	// OAuth and authentication
	app.Get("/api/v1/authenticate", endpoints.MakeGitlabAuthentication(store, oauthGitlab, clientId, oauthRedirectUri))
	app.Get("/oauth/redirect", endpoints.MakeOAuthAuthorize(store, oauthCallback, clientId, clientSecret, oauthDebug))

	go func() {
		stdLogger.Info("Application -> ListenTLS")
		if errTls := app.ListenTLS(":"+targetPort, "cert.pem", "key.pem"); errTls != nil {
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
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
	}
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
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

// Configure rdbms connection pool
func configureCnxPool(k *koanf.Koanf, zapLogger *zap.Logger) (*pgxpool.Pool, error) {

	dbUrl := k.String("app.pgUrl")
	poolMin := k.String("app.pgPoolMin")
	poolMax := k.String("app.pgPoolMax")

	zapLogger.Info("CnxPool -> Parse configuration")
	dbConfig, errDbCfg := pgxpool.ParseConfig(dbUrl)
	if errDbCfg != nil {
		panic(errDbCfg)
	}
	pgMin, errMin := strconv.Atoi(poolMin)
	if errMin != nil {
		panic(errMin)
	}
	pgMax, errMax := strconv.Atoi(poolMax)
	if errMax != nil {
		panic(errMax)
	}
	dbConfig.MinConns = int32(pgMin)
	dbConfig.MaxConns = int32(pgMax)

	zapLogger.Info("CnxPool -> Initialize pool")
	return pgxpool.NewWithConfig(context.Background(), dbConfig)
}
