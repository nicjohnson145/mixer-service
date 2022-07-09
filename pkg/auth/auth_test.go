package auth

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func setupDbAndRouter(t *testing.T) (*fiber.App, func()) {
	return commontest.SetupDbAndRouter(t, "auth.db", defineRoutes)
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

	T_RegisterUser(t, app, RegisterNewUserRequest{Username: "foo", Password: "bar"})

	for _, tc := range loginData {
		t.Run("login_cases_"+tc.name, func(t *testing.T) {
			if tc.expectedToken {
				status, resp := T_Login_ok(t, app, tc.input)
				require.Equal(t, tc.expectedCode, status)
				require.NotEmpty(t, resp.AccessToken)
				require.NotEmpty(t, resp.RefreshToken)
			} else {
				status, _ := T_Login_fail(t, app, tc.input)
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

	protectedRoute := func(c *fiber.Ctx, claims jwt.Claims) error {
		return c.SendString(respText)
	}
	app.Get("/some-protected-route", RequiresValidAccessToken(protectedRoute))

	tokenResetFunc := func() func() {
		currentTime := jwt.GetAccessTokenDuration()
		return func() {
			jwt.SetAccessTokenDuration(currentTime)
		}
	}
	resetTokenDuration := tokenResetFunc()
	defer resetTokenDuration()

	t_protectedRoute := func(t *testing.T, app *fiber.App, token string) (int, string) {
		req, err := http.NewRequest(http.MethodGet, "/some-protected-route", nil)
		require.NoError(t, err)
		req.Header.Set(jwt.AuthenticationHeader, token)

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
		req.Header.Set(jwt.AuthenticationHeader, refreshToken)

		resp, err := app.Test(req)
		require.NoError(t, err)

		defer resp.Body.Close()
		var r RefreshTokenResponse
		err = json.NewDecoder(resp.Body).Decode(&r)
		require.NoError(t, err)

		return resp.StatusCode, r
	}

	// Set the token expiry to something short, so we can get an "expired" token
	jwt.SetAccessTokenDuration(time.Duration(1 * time.Second))

	// Register a new user
	T_RegisterUser(t, app, RegisterNewUserRequest{Username: username, Password: password})

	// Successfully login
	status, loginResp := T_Login_ok(t, app, LoginRequest{Username: username, Password: password})
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

func TestChangePassword(t *testing.T) {
	app, cleanup := setupDbAndRouter(t)
	defer cleanup()

	// Register some new users
	T_RegisterUser(t, app, RegisterNewUserRequest{Username: "foo", Password: "bar"})
	T_RegisterUser(t, app, RegisterNewUserRequest{Username: "bar", Password: "baz"})

	// Login as each
	_, _ = T_Login_ok(t, app, LoginRequest{Username: "foo", Password: "bar"})
	_, _ = T_Login_ok(t, app, LoginRequest{Username: "bar", Password: "baz"})

	// Change foo's password
	T_ChangePassword_ok(t, app, ChangePasswordRequest{NewPassword: "barbar"}, commontest.AuthOpts{Username: commontest.Ptr("foo")})

	// Ensure logins are still kosher
	_, _ = T_Login_ok(t, app, LoginRequest{Username: "foo", Password: "barbar"})
	_, _ = T_Login_fail(t, app, LoginRequest{Username: "foo", Password: "bar"})
	_, _ = T_Login_ok(t, app, LoginRequest{Username: "bar", Password: "baz"})
}
