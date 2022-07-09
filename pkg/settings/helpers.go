package settings

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"

	"github.com/nicjohnson145/mixer-service/pkg/common"
)

func T_GetSettings(t *testing.T, app *fiber.App, o commontest.AuthOpts) (int, GetSettingsResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[any]{
		Method: http.MethodGet,
		Path:   common.SettingsV1,
		Auth:   &o,
	})
	return commontest.T_call_ok[GetSettingsResponse](t, app, req)
}

func T_WriteSettings(t *testing.T, app *fiber.App, b WriteSettingsRequest, o commontest.AuthOpts) (int, WriteSettingsResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[WriteSettingsRequest]{
		Method: http.MethodPut,
		Path:   common.SettingsV1,
		Body:   &b,
		Auth:   &o,
	})
	return commontest.T_call_ok[WriteSettingsResponse](t, app, req)
}

