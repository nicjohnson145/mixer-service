package common

import (
	"errors"
	"net/http"
	"os"
)

var ErrNotFound = errors.New("not found")

type HttpHandler func(http.ResponseWriter, *http.Request)

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
