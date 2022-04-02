package auth

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type User struct {
	Username string
	Password string
}

type ClaimsHttpHandler func(http.ResponseWriter, *http.Request, jwt.Claims)
type FiberClaimsHandler func(*fiber.Ctx, jwt.Claims) error

func Init(app *fiber.App, db *sql.DB) error {
	defineRoutes(app, db)
	return nil
}

func defineRoutes(app *fiber.App, db *sql.DB) {
	if common.DefaultedEnvVar("PROTECT_REGISTER_ENDPOINT", "false") == "true" {
		app.Post(common.AuthV1+"/register-user", RequiresValidAccessToken(registerNewUser(db)))
	} else {
		app.Post(common.AuthV1+"/register-user", noopProtection(registerNewUser(db)))
	}
	app.Post(common.AuthV1+"/login", login(db))
	app.Post(common.AuthV1+"/refresh", requiresValidRefreshToken(refresh()))
	app.Post(common.AuthV1+"/change-password", RequiresValidAccessToken(changePassword(db)))
}

func hashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func comparePasswords(hashedPw string, plainPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(plainPw))
	return err == nil
}

func noopProtection(handler FiberClaimsHandler) common.FiberHandler {
	return func(c *fiber.Ctx) error {
		return handler(c, jwt.Claims{})
	}
}

func requiresValidToken(handler FiberClaimsHandler, validationFunc func(string) (jwt.Claims, error)) common.FiberHandler {

	return func(c *fiber.Ctx) error {
		val := c.Get(jwt.AuthenticationHeader)
		if val == "" {
			return common.NewGenericUnauthorizedResponse("no authentication header")
		}

		claims, err := validationFunc(val)
		if err != nil {
			return common.NewGenericUnauthorizedResponse("invalid token")
		}

		return handler(c, claims)
	}
}

func RequiresValidAccessToken(handler FiberClaimsHandler) common.FiberHandler {
	return requiresValidToken(handler, jwt.ValidateAccessToken)
}

func requiresValidRefreshToken(handler FiberClaimsHandler) common.FiberHandler {
	return requiresValidToken(handler, jwt.ValidateRefreshToken)
}
