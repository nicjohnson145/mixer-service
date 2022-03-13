package common

import (
	"errors"
	"net/http"
)

var ErrNotFound = errors.New("not found")

type HttpHandler func(http.ResponseWriter, *http.Request)

const (
	ApiV1    = "/api/v1"
	AuthV1   = ApiV1 + "/auth"
	DrinksV1 = ApiV1 + "/drinks"
)
