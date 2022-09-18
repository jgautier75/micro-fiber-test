package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"micro-fiber-test/pkg/dao/impl"
	"micro-fiber-test/pkg/endpoints"
	svcImpl "micro-fiber-test/pkg/service/impl"
)

func main() {

	// Setup service & dao
	orgDao := impl.OrgDao{}
	orgSvc := svcImpl.NewOrgService(&orgDao)

	// Load config file
	var k = koanf.New(".")
	err := k.Load(file.Provider("config/config.yaml"), yaml.Parser())
	if err != nil {
		fmt.Printf("Error loading confguration file [%v]", err)
	}
	targetPort := k.String("http.server.port")
	defaultTenantId := k.Int64("app.tenant")
	dbUrl := k.String("app.pgUrl")

	app := fiber.New()

	app.Post("/api/v1/organizations", endpoints.MakeOrgCreateEndpoint(dbUrl, defaultTenantId, orgSvc))

	app.ListenTLS(":"+targetPort, "cert.pem", "key.pem")
}
