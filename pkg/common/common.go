package common

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var ErrNotFound = errors.New("not found")

type HttpHandler func(http.ResponseWriter, *http.Request)
type FiberHandler func(*fiber.Ctx) error

const (
	ApiV1      = "/api/v1"
	AuthV1     = ApiV1 + "/auth"
	DrinksV1   = ApiV1 + "/drinks"
	HealthV1   = ApiV1 + "/health"
	SettingsV1 = ApiV1 + "/settings"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Immutable:    true,
		ErrorHandler: ErrorHandler,
	})
	app.Use(logger.New())
	return app
}

func DefaultedEnvVar(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} else {
		return defaultVal
	}
}

type ErrorResponse struct {
	Msg     string
	Err     error
	Context string
	Status  int
}

func (e ErrorResponse) Error() string {
	return e.Err.Error()
}

type OutboundErrResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

func NewInternalServerErrorResp(context string, e error) ErrorResponse {
	return ErrorResponse{
		Msg:     "internal server error",
		Err:     e,
		Context: context,
		Status:  fiber.StatusInternalServerError,
	}
}

func NewGenericUnauthorizedResponse(context string) ErrorResponse {
	return ErrorResponse{
		Msg:     "unauthorized",
		Err:     fmt.Errorf("unauthorized"),
		Context: context,
		Status:  fiber.StatusUnauthorized,
	}
}

func NewGenericNotFoundResponse(context string) ErrorResponse {
	return ErrorResponse{
		Msg:     "not found",
		Err:     fmt.Errorf("not found"),
		Context: context,
		Status:  fiber.StatusNotFound,
	}
}

func NewBadRequestResponse(e error) ErrorResponse {
	return ErrorResponse{
		Msg:     "bad request",
		Err:     e,
		Context: "",
		Status:  fiber.StatusBadRequest,
	}
}

func ErrorHandler(c *fiber.Ctx, e error) error {
	er, ok := e.(ErrorResponse)
	if !ok {
		return fiber.DefaultErrorHandler(c, e)
	}

	status := fiber.StatusInternalServerError
	if er.Status != 0 {
		status = er.Status
	}

	logMsg := er.Err.Error()
	if er.Context != "" {
		logMsg = fmt.Sprintf("(%v) %v", er.Context, er.Err.Error())
	}
	log.Error(logMsg)

	return c.Status(status).JSON(OutboundErrResponse{
		Error:   er.Msg,
		Success: false,
	})
}
