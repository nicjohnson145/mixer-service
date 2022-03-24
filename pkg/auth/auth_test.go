package auth

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func setupDbAndRouter(t *testing.T) (*fiber.App, func()) {
	return commontest.SetupDbAndRouter(t, "auth.db", defineRoutes)
}

func t_registerUser(t *testing.T, app *fiber.App, req RegisterNewUserRequest) {
	body, err := json.Marshal(req)
	require.NoError(t, err)

	registerReq, err := http.NewRequest("POST", common.AuthV1+"/register-user", strings.NewReader(string(body)))
	commontest.SetJsonHeader(registerReq)
	require.NoError(t, err)

	resp, err := app.Test(registerReq)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func t_login(t *testing.T, app *fiber.App, req LoginRequest) *http.Response {
	body, err := json.Marshal(req)
	require.NoError(t, err)
	loginReq, err := http.NewRequest("POST", common.AuthV1+"/login", strings.NewReader(string(body)))
	require.NoError(t, err)
	commontest.SetJsonHeader(loginReq)
	resp, err := app.Test(loginReq)
	require.NoError(t, err)

	return resp
}

func t_login_ok(t *testing.T, app *fiber.App, req LoginRequest) (int, LoginResponse) {
	resp := t_login(t, app, req)
	commontest.RequireOkStatus(t, resp)

	var r LoginResponse
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(&r)
	require.NoError(t, err)
	return resp.StatusCode, r
}

func t_login_fail(t *testing.T, app *fiber.App, req LoginRequest) (int, common.OutboundErrResponse) {
	resp := t_login(t, app, req)
	commontest.RequireNotOkStatus(t, resp)

	var r common.OutboundErrResponse
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(&r)
	require.NoError(t, err)
	return resp.StatusCode, r
}

func TestRegisterLogin(t *testing.T) {
	loginData := []struct {
		name          string
		input         LoginRequest
		expectedCode  int
		expectedToken bool
	}{
		{
			name:          "missing_user",
			input:         LoginRequest{Username: "not_a_user", Password: "bar"},
			expectedCode:  http.StatusUnauthorized,
			expectedToken: false,
		},
		{
			name:          "incorrect_password",
			input:         LoginRequest{Username: "foo", Password: "wrong_password"},
			expectedCode:  http.StatusUnauthorized,
			expectedToken: false,
		},
		{
			name:          "valid_login",
			input:         LoginRequest{Username: "foo", Password: "bar"},
			expectedCode:  http.StatusOK,
			expectedToken: true,
		},
	}

	app, cleanup := setupDbAndRouter(t)
	defer cleanup()

	t_registerUser(t, app, RegisterNewUserRequest{Username: "foo", Password: "bar"})

	for _, tc := range loginData {
		t.Run("login_cases_"+tc.name, func(t *testing.T) {
			if tc.expectedToken {
				status, resp := t_login_ok(t, app, tc.input)
				require.Equal(t, tc.expectedCode, status)
				require.NotEmpty(t, resp.AccessToken)
				require.NotEmpty(t, resp.RefreshToken)
			} else {
				status, _ := t_login_fail(t, app, tc.input)
				require.Equal(t, status, tc.expectedCode)
			}
		})
	}
}

func TestRefresh(t *testing.T) {
	const (
		respText = "yer authorized"
		username = "foo"
		password = "bar"
	)

	app, cleanup := setupDbAndRouter(t)
	defer cleanup()

	protectedRoute := func(c *fiber.Ctx, claims Claims) error {
		return c.SendString(respText)
	}
	app.Get("/some-protected-route", RequiresValidAccessToken(protectedRoute))

	tokenResetFunc := func() func() {
		currentTime := accessTokenDuration
		return func() {
			accessTokenDuration = currentTime
		}
	}
	resetTokenDuration := tokenResetFunc()
	defer resetTokenDuration()

	t_protectedRoute := func(t *testing.T, app *fiber.App, token string) (int, string) {
		req, err := http.NewRequest(http.MethodGet, "/some-protected-route", nil)
		require.NoError(t, err)
		req.Header.Set(AuthenticationHeader, token)

		resp, err := app.Test(req)
		require.NoError(t, err)

		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		return resp.StatusCode, string(bytes)
	}

	validProtectedRoute := func(t *testing.T, app *fiber.App, token string) {
		status, body := t_protectedRoute(t, app, token)
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, body, respText)
	}

	unauthorizedProtectedRoute := func(t *testing.T, app *fiber.App, token string) {
		status, _ := t_protectedRoute(t, app, token)
		require.Equal(t, http.StatusUnauthorized, status)
	}

	t_refresh := func(t *testing.T, app *fiber.App, refreshToken string) (int, RefreshTokenResponse) {
		req, err := http.NewRequest(http.MethodPost, common.AuthV1+"/refresh", nil)
		require.NoError(t, err)
		req.Header.Set(AuthenticationHeader, refreshToken)

		resp, err := app.Test(req)
		require.NoError(t, err)

		defer resp.Body.Close()
		var r RefreshTokenResponse
		err = json.NewDecoder(resp.Body).Decode(&r)
		require.NoError(t, err)

		return resp.StatusCode, r
	}

	// Set the token expiry to something short, so we can get an "expired" token
	accessTokenDuration = time.Duration(1 * time.Second)

	// Register a new user
	t_registerUser(t, app, RegisterNewUserRequest{Username: username, Password: password})

	// Successfully login
	status, loginResp := t_login_ok(t, app, LoginRequest{Username: username, Password: password})
	require.Equal(t, http.StatusOK, status)

	// Hit a protected route to prove our token is valid
	validProtectedRoute(t, app, loginResp.AccessToken)

	// Wait a bit for the token to expire (not great I know)
	time.Sleep(2 * time.Second)

	// Token should be invalid by now, but our refresh token shouldnt
	unauthorizedProtectedRoute(t, app, loginResp.AccessToken)

	// Refresh the token
	status, refreshResponse := t_refresh(t, app, loginResp.RefreshToken)
	require.Equal(t, http.StatusOK, status)

	// Should be able to use the new token to hit the protected route
	validProtectedRoute(t, app, refreshResponse.AccessToken)
}
