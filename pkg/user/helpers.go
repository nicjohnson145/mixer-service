package user

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"

	"github.com/nicjohnson145/mixer-service/pkg/common"
)

func T_GetPublicUsers(t *testing.T, app *fiber.App, o commontest.AuthOpts) (int, GetPublicUsersResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[any]{
		Method: http.MethodGet,
		Path:   common.UserV1,
		Auth:   &o,
	})
	return commontest.T_call_ok[GetPublicUsersResponse](t, app, req)
}
