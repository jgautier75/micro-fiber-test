package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
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
)

func main() {

	// Load config file
	var k = koanf.New(".")
	err := k.Load(file.Provider("config/config.yaml"), yaml.Parser())
	if err != nil {
		fmt.Printf("Error loading confguration file [%v]", err)
	}
	targetPort := k.String("http.server.port")
	defaultTenantId := k.Int64("app.tenant")
	dbUrl := k.String("app.pgUrl")
	logCfg := k.String("app.logFile")

	// Setup service & dao
	orgDao := impl.NewOrgDao(dbUrl)
	sectorDao := impl.NewSectorDao(dbUrl)
	userDao := impl.NewUserDao(dbUrl)
	orgSvc := svcImpl.NewOrgService(orgDao, sectorDao)
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

	fConfig := fiber.Config{
		CaseSensitive:     true,
		StrictRouting:     true,
		EnablePrintRoutes: false,
		UnescapePath:      true,
		ErrorHandler:      defErrorHandler,
	}

	app := fiber.New(fConfig)

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
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

	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	app.Use(logging.New(zapLogger))

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

	errTls := app.ListenTLS(":"+targetPort, "cert.pem", "key.pem")
	if errTls != nil {
		fmt.Printf("ListenTLS error [%s]", errTls)
	}
}
