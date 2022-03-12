package common

import (
	"errors"
	"net/http"
)

var ErrNotFound = errors.New("not found")

type HttpHandler func(http.ResponseWriter, *http.Request)
