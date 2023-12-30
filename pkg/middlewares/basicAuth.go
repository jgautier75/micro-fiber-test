package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func NewBasicAuthConfig(refUserName string, refPass string) basicauth.Config {
	return basicauth.Config{
		Users: map[string]string{
			refUserName: refPass,
		},
		Realm: "Forbidden",
		Authorizer: func(user string, pass string) bool {
			if user == refUserName && pass == refPass {
				return true
			}
			return false
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Redirect("/unauthorized.html", 403)
		},
		ContextUsername: "_user",
		ContextPassword: "_pass",
	}

}
