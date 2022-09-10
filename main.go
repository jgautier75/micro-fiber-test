package main

import (
	"github.com/gofiber/fiber/v2"
	"micro-fiber-test/pkg/dao/impl"
	"micro-fiber-test/pkg/endpoints"
	svcImpl "micro-fiber-test/pkg/service/impl"
)

const (
	defaultTenantId = 1
	rdbmsUrl        = "postgres://unicorn_user:magical_password@localhost:5432/rainbow_database"
)

func main() {

	// Setup service & dao
	orgDao := impl.OrgDao{}
	orgSvc := svcImpl.NewOrgService(&orgDao)

	app := fiber.New()

	app.Post("/api/v1/organizations", endpoints.MakeOrgCreateEndpoint(rdbmsUrl, defaultTenantId, orgSvc))

	app.Listen(":8080")
}
