package common

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
)

var ErrNotFound = errors.New("not found")

type HttpHandler func(http.ResponseWriter, *http.Request)
type FiberHandler func(*fiber.Ctx) error

const (
	ApiV1    = "/api/v1"
	AuthV1   = ApiV1 + "/auth"
	DrinksV1 = ApiV1 + "/drinks"
	HealthV1 = ApiV1 + "/health"
)

func DefaultedEnvVar(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} else {
		return defaultVal
	}
}
