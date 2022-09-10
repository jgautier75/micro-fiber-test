package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"micro-fiber-test/pkg/commons"
	"micro-fiber-test/pkg/contracts"
	"micro-fiber-test/pkg/dao/impl"
	"micro-fiber-test/pkg/model"
	svcImpl "micro-fiber-test/pkg/service/impl"
	"net/http"
)

const rdbmsUrl = "postgres://unicorn_user:magical_password@localhost:5432/rainbow_database"

func main() {

	// Setup service & dao
	orgDao := impl.OrgDao{}
	orgSvc := svcImpl.NewOrgService(&orgDao)

	app := fiber.New()

	app.Post("/api/v1/organizations", func(ctx *fiber.Ctx) error {
		org := model.Organization{}
		if err := ctx.BodyParser(org); err != nil {
			fmt.Println("error = ", err)
			return ctx.SendStatus(200)
		}
		id, err := orgSvc.Create(rdbmsUrl, &org)
		if err != nil {
			ctx.SendStatus(http.StatusInternalServerError)
			apiErr := commons.ApiError{
				Code:    http.StatusInternalServerError,
				Kind:    string(commons.ErrorTypeTechnical),
				Message: err.Error(),
			}
			return ctx.JSON(apiErr)
		} else {
			idResponse := contracts.IdResponse{ID: id}
			return ctx.JSON(idResponse)
		}

	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":8080")
}
