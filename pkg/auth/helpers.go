package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"
	"net/http"
	"testing"
)

func T_RegisterUser(t *testing.T, app *fiber.App, r RegisterNewUserRequest) (int, RegisterNewUserResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[RegisterNewUserRequest]{
		Method: http.MethodPost,
		Path:   common.AuthV1 + "/register-user",
		Body:   &r,
	})
	return commontest.T_call_ok[RegisterNewUserResponse](t, app, req)
}

func T_Login_ok(t *testing.T, app *fiber.App, r LoginRequest) (int, LoginResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[LoginRequest]{
		Method: http.MethodPost,
		Path:   common.AuthV1 + "/login",
		Body:   &r,
	})
	return commontest.T_call_ok[LoginResponse](t, app, req)
}

func T_Login_fail(t *testing.T, app *fiber.App, r LoginRequest) (int, common.OutboundErrResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[LoginRequest]{
		Method: http.MethodPost,
		Path:   common.AuthV1 + "/login",
		Body:   &r,
	})
	return commontest.T_call_fail(t, app, req)
}

func T_ChangePassword_ok(t *testing.T, app *fiber.App, r ChangePasswordRequest, o commontest.AuthOpts) (int, ChangePasswordResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[ChangePasswordRequest]{
		Method: http.MethodPost,
		Path:   common.AuthV1 + "/change-password",
		Auth:   &o,
		Body:   &r,
	})
	return commontest.T_call_ok[ChangePasswordResponse](t, app, req)
}
