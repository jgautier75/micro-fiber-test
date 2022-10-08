package endpoints

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"io"
	"micro-fiber-test/pkg/commons"
	"micro-fiber-test/pkg/contracts"
	"micro-fiber-test/pkg/model"
	"net/http"
	"net/url"
	"strings"
)

const oAuthState = "oauthstate"

func MakeOAuthAuthorize(store *session.Store, oauthCallback string, oAuthClientId string, oauthClientSecret string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		httpClient := http.Client{}
		code := ctx.Query("code")
		reqState := ctx.Query("state")

		// Decode state from query
		decState, errDecode := url.QueryUnescape(reqState)
		if errDecode != nil {
			contracts.ConvertToInternalError(errDecode)
		}

		// Get state from session and decode
		httpSession, errSession := store.Get(ctx)
		if errSession != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := contracts.ConvertToInternalError(errSession)
			return ctx.JSON(apiError)
		}
		sessionOAuth := httpSession.Get(oAuthState)
		decodedSessionOAuth, errDecodeSessionOAuth := url.QueryUnescape(sessionOAuth.(string))
		if errDecodeSessionOAuth != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := contracts.ConvertToInternalError(errDecodeSessionOAuth)
			return ctx.JSON(apiError)
		}

		// Compare http request state and state from session
		if decodedSessionOAuth != decState {
			apiError := contracts.ConvertToFunctionalError(errors.New(commons.OAuthStateMismatch), fiber.StatusConflict)
			_ = ctx.SendStatus(fiber.StatusConflict)
			return ctx.JSON(apiError)
		}

		// Delete from session
		httpSession.Delete(oAuthState)

		reqURL := fmt.Sprintf(oauthCallback, oAuthClientId, oauthClientSecret, code)
		req, errOauth := http.NewRequest(http.MethodPost, reqURL, nil)
		if errOauth != nil {
			apiError := contracts.ConvertToFunctionalError(errOauth, fiber.StatusConflict)
			_ = ctx.SendStatus(fiber.StatusConflict)
			return ctx.JSON(apiError)
		}
		req.Header.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

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

func MakeGitlabAuthentication(store *session.Store, oauthGitlab string, clientId string, redirectUri string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var s strings.Builder
		state, errState := generateState(28)
		if errState != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := contracts.ConvertToInternalError(errState)
			return ctx.JSON(apiError)
		}

		httpSession, errSession := store.Get(ctx)
		if errSession != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := contracts.ConvertToInternalError(errSession)
			return ctx.JSON(apiError)
		}
		httpSession.Set(oAuthState, state)
		_ = httpSession.Save()

		s.WriteString(oauthGitlab)
		s.WriteString("?client_id=")
		s.WriteString(clientId)
		s.WriteString("&response_type=code")
		s.WriteString("&state=")
		s.WriteString(state)
		s.WriteString("&scope=user")
		s.WriteString("&redirect_uri=")
		s.WriteString(redirectUri)
		return ctx.Redirect(s.String())
	}
}

// Generate random state of size n
func generateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(data)
	return state, nil
}
