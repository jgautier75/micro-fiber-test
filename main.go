package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"micro-fiber-test/pkg/contracts"
	"micro-fiber-test/pkg/dao/impl"
	"micro-fiber-test/pkg/endpoints"
	"micro-fiber-test/pkg/logging"
	svcImpl "micro-fiber-test/pkg/service/impl"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {

	// Load config file
	var k = koanf.New(".")
	errLoadCfg := k.Load(file.Provider("config/config.yaml"), yaml.Parser())
	if errLoadCfg != nil {
		panic(errLoadCfg)
	}
	targetPort := k.String("http.server.port")
	defaultTenantId := k.Int64("app.tenant")
	logCfg := k.String("app.logFile")
	clientId := k.String("app.oauthClientId")
	clientSecret := k.String("app.oauthClientSecret")
	oauthCallback := k.String("app.oauthCallback")
	oauthRedirectUri := k.String("app.oauthRedirectUri")
	oauthGitlab := k.String("app.oauthGitlab")
	dbUrl := k.String("app.pgUrl")

	zapLogger := configureLogger(logCfg)

	zapLogger.Info("OAuth2",
		zap.Field{Key: "clientId", Type: zapcore.StringType, String: clientId},
		zap.Field{Key: "clientSecret", Type: zapcore.StringType, String: clientSecret},
	)

	dbPool, poolErr := configureCnxPool(k, zapLogger)
	if poolErr != nil {
		panic(poolErr)
	}

	zapLogger.Info("Dao & Services -> Setup & inject")
	orgDao := impl.NewOrgDao(dbPool)
	sectorDao := impl.NewSectorDao(dbPool)
	userDao := impl.NewUserDao(dbPool)
	orgSvc := svcImpl.NewOrgService(dbPool, orgDao, sectorDao)
	sectorSvc := svcImpl.NewSectorService(sectorDao)
	userSvc := svcImpl.NewUserService(userDao)

	var defErrorHandler = func(c *fiber.Ctx, err error) error {
		var e *fiber.Error
		code := fiber.StatusInternalServerError
		if errors.As(err, &e) {
			code = e.Code
			if code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError {
				apiError := contracts.ConvertToFunctionalError(err, code)
				return c.Status(code).JSON(apiError)
			} else {
				apiError := contracts.ConvertToInternalError(err)
				return c.Status(code).JSON(apiError)
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(contracts.ConvertToInternalError(err))
	}

	defCfg := session.ConfigDefault
	store := session.New(defCfg)

	fConfig := fiber.Config{
		CaseSensitive:     true,
		StrictRouting:     true,
		EnablePrintRoutes: false,
		UnescapePath:      true,
		ErrorHandler:      defErrorHandler,
	}

	zapLogger.Info("Application -> Setup")
	app := fiber.New(fConfig)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	serverShutdown := make(chan struct{})
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
		serverShutdown <- struct{}{}
	}()

	app.Use(logging.New(zapLogger))

	app.Static("/", "./static")

	// Organizations
	app.Post("/api/v1/organizations", endpoints.MakeOrgCreateEndpoint(dbUrl, defaultTenantId, orgSvc))
	app.Put("/api/v1/organizations/:orgCode", endpoints.MakeOrgUpdateEndpoint(defaultTenantId, orgSvc))
	app.Delete("/api/v1/organizations/:orgCode", endpoints.MakeOrgDeleteEndpoint(defaultTenantId, orgSvc))
	app.Get("/api/v1/organizations/:orgCode", endpoints.MakeOrgFindByCodeEndpoint(defaultTenantId, orgSvc))
	app.Get("/api/v1/organizations", endpoints.MakeOrgFindAll(defaultTenantId, orgSvc))

	// Sectors
	app.Get("/api/v1/organizations/:orgCode/sectors", endpoints.MakeSectorsFindByOrga(defaultTenantId, orgSvc, sectorSvc))
	app.Post("/api/v1/organizations/:orgCode/sectors", endpoints.MakeSectorCreateEndpoint(defaultTenantId, orgSvc, sectorSvc))
	app.Delete("/api/v1/organizations/:orgCode/sectors/:sectorCode", endpoints.MakeSectorDeleteEndpoint(defaultTenantId, orgSvc, sectorSvc))
	app.Put("/api/v1/organizations/:orgCode/sectors/:sectorCode", endpoints.MakeSectorUpdateEndpoint(defaultTenantId, orgSvc, sectorSvc))

	// Users
	app.Post("/api/v1/organizations/:orgCode/users", endpoints.MakeUserCreateEndpoint(defaultTenantId, userSvc, orgSvc))
	app.Get("/api/v1/organizations/:orgCode/users", endpoints.MakeUserSearchFilter(defaultTenantId, userSvc, orgSvc))
	app.Get("/api/v1/organizations/:orgCode/users/:userId", endpoints.MakeUserFindByCode(defaultTenantId, userSvc, orgSvc))
	app.Put("/api/v1/organizations/:orgCode/users/:userId", endpoints.MakeUserUpdate(defaultTenantId, userSvc, orgSvc))
	app.Delete("/api/v1/organizations/:orgCode/users/:userId", endpoints.MakeUserDelete(defaultTenantId, userSvc, orgSvc))

	// OAuth and authentication
	app.Get("/api/v1/authenticate", endpoints.MakeGitlabAuthentication(store, oauthGitlab, clientId, oauthRedirectUri))
	app.Get("/oauth/redirect", endpoints.MakeOAuthAuthorize(store, oauthCallback, clientId, clientSecret))

	zapLogger.Info("Application -> ListenTLS")
	errTls := app.ListenTLS(":"+targetPort, "cert.pem", "key.pem")
	if errTls != nil {
		fmt.Printf("ListenTLS error [%s]", errTls)
	}

	zapLogger.Info("Closing connections pool")
	dbPool.Close()

}

func configureLogger(logCfg string) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
	}
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile(logCfg, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return zapLogger
}

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
