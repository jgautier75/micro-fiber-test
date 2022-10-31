package endpoints

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"io"
	"micro-fiber-test/pkg/dto/commons"
	"micro-fiber-test/pkg/exceptions"
	"micro-fiber-test/pkg/model"
	"net/http"
	"net/url"
	"strings"
)

const (
	oAuthState = "oauthstate"
	cVerifier  = "cverifier"
)

func MakeOAuthAuthorize(store *session.Store, oauthCallback string, oAuthClientId string, oauthClientSecret string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		httpClient := http.Client{}
		code := ctx.Query("code")
		reqState := ctx.Query("state")

		// Decode state from query
		decState, errDecode := url.QueryUnescape(reqState)
		if errDecode != nil {
			exceptions.ConvertToInternalError(errDecode)
		}

		// Get state from session and decode
		httpSession, errSession := store.Get(ctx)
		if errSession != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := exceptions.ConvertToInternalError(errSession)
			return ctx.JSON(apiError)
		}
		sessionOAuth := httpSession.Get(oAuthState)
		decodedSessionOAuth, errDecodeSessionOAuth := url.QueryUnescape(sessionOAuth.(string))
		if errDecodeSessionOAuth != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := exceptions.ConvertToInternalError(errDecodeSessionOAuth)
			return ctx.JSON(apiError)
		}

		// Compare http request state and state from session
		if decodedSessionOAuth != decState {
			apiError := exceptions.ConvertToFunctionalError(errors.New(commons.OAuthStateMismatch), fiber.StatusConflict)
			_ = ctx.SendStatus(fiber.StatusConflict)
			return ctx.JSON(apiError)
		}

		codeVerifier := httpSession.Get(cVerifier)
		reqURL := fmt.Sprintf(oauthCallback, oAuthClientId, oauthClientSecret, code, codeVerifier)

		// Delete from session
		httpSession.Delete(oAuthState)
		httpSession.Delete(cVerifier)

		req, errOauth := http.NewRequest(http.MethodPost, reqURL, nil)
		if errOauth != nil {
			apiError := exceptions.ConvertToFunctionalError(errOauth, fiber.StatusConflict)
			_ = ctx.SendStatus(fiber.StatusConflict)
			return ctx.JSON(apiError)
		}
		req.Header.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

		res, errHttp := httpClient.Do(req)
		if errHttp != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := exceptions.ConvertToInternalError(errHttp)
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
			apiError := exceptions.ConvertToInternalError(errDecode)
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
			apiError := exceptions.ConvertToInternalError(errState)
			return ctx.JSON(apiError)
		}

		httpSession, errSession := store.Get(ctx)
		if errSession != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := exceptions.ConvertToInternalError(errSession)
			return ctx.JSON(apiError)
		}
		httpSession.Set(oAuthState, state)

		// Generate code verifier
		buf, errRnd := randomBytes(32)
		if errRnd != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := exceptions.ConvertToInternalError(errState)
			return ctx.JSON(apiError)
		}
		encodedRandom := encode(buf)
		httpSession.Set(cVerifier, encodedRandom)
		_ = httpSession.Save()

		h := sha256.New()
		_, errSha := h.Write([]byte(encodedRandom))
		shaChallenge := encode(h.Sum(nil))

		if errSha != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiError := exceptions.ConvertToInternalError(errSha)
			return ctx.JSON(apiError)
		}

		s.WriteString(oauthGitlab)
		s.WriteString("?client_id=")
		s.WriteString(clientId)
		s.WriteString("&response_type=code")
		s.WriteString("&code_challenge=")
		s.WriteString(shaChallenge)
		s.WriteString("&code_challenge_method=S256")
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

// Encode code verifier according to protect against CSRF attacks
func encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)
	return encoded
}

// Generate a random string of specified length
func randomBytes(length int) ([]byte, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const csLen = byte(len(charset))
	output := make([]byte, 0, length)
	for {
		buf := make([]byte, length)
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			return nil, fmt.Errorf("failed to read random bytes: %v", err)
		}
		for _, b := range buf {
			// Avoid bias by using a value range that's a multiple of 62
			if b < (csLen * 4) {
				output = append(output, charset[b%csLen])

				if len(output) == length {
					return output, nil
				}
			}
		}
	}
}