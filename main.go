package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"micro-fiber-test/pkg/commons"
	"micro-fiber-test/pkg/contracts"
	"micro-fiber-test/pkg/dao/impl"
	"micro-fiber-test/pkg/model"
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

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Page %s not found", c.OriginalURL())
		apiErr := commons.ApiError{
			Code:    fiber.StatusNotFound,
			Kind:    string(commons.ErrorTypeFunctional),
			Message: msg,
		}
		c.SendStatus(404)
		return c.JSON(apiErr)
	})

	app.Post("/api/v1/organizations", func(ctx *fiber.Ctx) error {
		payload := struct {
			Code   string `json:"code"`
			Label  string `json:"label"`
			Kind   string `json:"type"`
			Status int    `json:"status"`
		}{}
		if err := ctx.BodyParser(&payload); err != nil {
			fmt.Println("error = ", err)
			return ctx.SendStatus(200)
		}
		org := model.Organization{}
		org.SetTenantId(defaultTenantId)
		org.SetCode(payload.Code)
		org.SetLabel(payload.Label)
		org.SetType(model.OrganizationType(payload.Kind))
		org.SetStatus(model.OrganizationStatus(payload.Status))
		id, err := orgSvc.Create(rdbmsUrl, &org)
		if err != nil {
			ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := commons.ApiError{
				Code:    fiber.StatusInternalServerError,
				Kind:    string(commons.ErrorTypeTechnical),
				Message: err.Error(),
			}
			return ctx.JSON(apiErr)
		} else {
			ctx.SendStatus(fiber.StatusCreated)
			idResponse := contracts.IdResponse{ID: id}
			return ctx.JSON(idResponse)
		}

	})

	app.Listen(":8080")
}
