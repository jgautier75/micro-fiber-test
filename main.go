package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"micro-fiber-test/pkg/contracts"
	"micro-fiber-test/pkg/dao/impl"
	"micro-fiber-test/pkg/endpoints"
	svcImpl "micro-fiber-test/pkg/service/impl"
)

func main() {

	// Setup service & dao
	orgDao := impl.NewOrgDao()
	sectorDao := impl.NewSectorDao()
	userDao := impl.NewUserDao()
	orgSvc := svcImpl.NewOrgService(orgDao, sectorDao)
	sectorSvc := svcImpl.NewSectorService(sectorDao)
	userSvc := svcImpl.NewUserService(userDao)

	// Load config file
	var k = koanf.New(".")
	err := k.Load(file.Provider("config/config.yaml"), yaml.Parser())
	if err != nil {
		fmt.Printf("Error loading confguration file [%v]", err)
	}
	targetPort := k.String("http.server.port")
	defaultTenantId := k.Int64("app.tenant")
	dbUrl := k.String("app.pgUrl")

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

	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05-0700",
		TimeZone:   "UTC",
		Format:     "[${time}] - [${ip}]:${port} ${status} - ${method} - ${path}\n<<<<<<<<<< Request\n${reqHeaders}\n${body}\n>>>>>>>>>> Response\n${protocol}:${status}\nBody:${resBody}\n",
	}))

	// Organizations
	app.Post("/api/v1/organizations", endpoints.MakeOrgCreateEndpoint(dbUrl, defaultTenantId, orgSvc))
	app.Put("/api/v1/organizations/:orgCode", endpoints.MakeOrgUpdateEndpoint(dbUrl, defaultTenantId, orgSvc))
	app.Delete("/api/v1/organizations/:orgCode", endpoints.MakeOrgDeleteEndpoint(dbUrl, defaultTenantId, orgSvc))
	app.Get("/api/v1/organizations/:orgCode", endpoints.MakeOrgFindByCodeEndpoint(dbUrl, defaultTenantId, orgSvc))
	app.Get("/api/v1/organizations", endpoints.MakeOrgFindAll(dbUrl, defaultTenantId, orgSvc))

	// Sectors
	app.Get("/api/v1/organizations/:orgCode/sectors", endpoints.MakeSectorsFindByOrga(dbUrl, defaultTenantId, orgSvc, sectorSvc))
	app.Post("/api/v1/organizations/:orgCode/sectors", endpoints.MakeSectorCreateEndpoint(dbUrl, defaultTenantId, orgSvc, sectorSvc))
	app.Delete("/api/v1/organizations/:orgCode/sectors/:sectorCode", endpoints.MakeSectorDeleteEndpoint(dbUrl, defaultTenantId, orgSvc, sectorSvc))

	// Users
	app.Post("/api/v1/organizations/:orgCode/users", endpoints.MakeUserCreateEndpoint(dbUrl, defaultTenantId, userSvc, orgSvc))
	app.Get("/api/v1/organizations/:orgCode/users", endpoints.MakeUserSearchFilter(dbUrl, defaultTenantId, userSvc, orgSvc))

	errTls := app.ListenTLS(":"+targetPort, "cert.pem", "key.pem")
	if errTls != nil {
		fmt.Printf("ListenTLS error [%s]", errTls)
	}
}
