package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"micro-fiber-test/pkg/contracts"
	"micro-fiber-test/pkg/model"
	"net/http"
)

func MakeOAuthAuthorize(oauthCallback string, oAuthClientId string, oauthClientSecret string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		httpClient := http.Client{}
		code := ctx.Query("code")
		reqURL := fmt.Sprintf(oauthCallback, oAuthClientId, oauthClientSecret, code)
		req, errOauth := http.NewRequest(http.MethodPost, reqURL, nil)
		if errOauth != nil {
			apiError := contracts.ConvertToFunctionalError(errOauth, fiber.StatusConflict)
			_ = ctx.SendStatus(fiber.StatusConflict)
			return ctx.JSON(apiError)
		}
		req.Header.Set("accept", "application/json")

		// Send out the HTTP request
		res, errHttp := httpClient.Do(req)
		if errHttp != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := contracts.ConvertToInternalError(errHttp)
			return ctx.JSON(apiError)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(res.Body)
		var t model.OAuthAccessResponse
		if errDecode := json.NewDecoder(res.Body).Decode(&t); errDecode != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := contracts.ConvertToInternalError(errDecode)
			return ctx.JSON(apiError)
		}
		return ctx.Redirect("/welcome.html?access_token=" + t.AccessToken)
	}
}
