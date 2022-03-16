package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setupDbAndRouter(t *testing.T) (*mux.Router, func()) {
	return commontest.SetupDbAndRouter(t, "auth.db", defineRoutes)
}

func t_registerUser(t *testing.T, router *mux.Router, req RegisterNewUserRequest) {
	body, err := json.Marshal(req)
	require.NoError(t, err)

	registerReq, err := http.NewRequest("POST", common.AuthV1+"/register-user", strings.NewReader(string(body)))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, registerReq)
	require.Equal(t, http.StatusOK, rr.Result().StatusCode)
}

func t_login(t *testing.T, router *mux.Router, req LoginRequest) (int, LoginResponse) {
	body, err := json.Marshal(req)
	require.NoError(t, err)
	loginReq, err := http.NewRequest("POST", common.AuthV1+"/login", strings.NewReader(string(body)))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, loginReq)

	defer rr.Result().Body.Close()
	var resp LoginResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&resp)
	require.NoError(t, err)

	return rr.Result().StatusCode, resp
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

	router, cleanup := setupDbAndRouter(t)
	defer cleanup()

	t_registerUser(t, router, RegisterNewUserRequest{Username: "foo", Password: "bar"})

	for _, tc := range loginData {
		t.Run("login_cases_"+tc.name, func(t *testing.T) {
			status, resp := t_login(t, router, tc.input)
			require.Equal(t, tc.expectedCode, status)

			if tc.expectedToken {
				require.NotEmpty(t, resp.AccessToken)
				require.NotEmpty(t, resp.RefreshToken)
			} else {
				require.Empty(t, resp.AccessToken)
				require.Empty(t, resp.RefreshToken)
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

	router, cleanup := setupDbAndRouter(t)
	defer cleanup()

	protectedRoute := func(w http.ResponseWriter, req *http.Request, claims Claims) {
		fmt.Fprint(w, respText)
	}
	router.HandleFunc("/some-protected-route", RequiresValidAccessToken(protectedRoute)).Methods(http.MethodGet)

	tokenResetFunc := func() func() {
		currentTime := accessTokenDuration
		return func() {
			accessTokenDuration = currentTime
		}
	}
	resetTokenDuration := tokenResetFunc()
	defer resetTokenDuration()

	t_protectedRoute := func(t *testing.T, router *mux.Router, token string) (int, string) {
		req, err := http.NewRequest(http.MethodGet, "/some-protected-route", nil)
		require.NoError(t, err)
		req.Header.Set(AuthenticationHeader, token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		defer rr.Result().Body.Close()
		bytes, err := ioutil.ReadAll(rr.Result().Body)
		require.NoError(t, err)
		return rr.Result().StatusCode, string(bytes)
	}

	validProtectedRoute := func(t *testing.T, router *mux.Router, token string) {
		status, body := t_protectedRoute(t, router, token)
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, body, respText)
	}

	unauthorizedProtectedRoute := func(t *testing.T, router *mux.Router, token string) {
		status, _ := t_protectedRoute(t, router, token)
		require.Equal(t, http.StatusUnauthorized, status)
	}

	t_refresh := func(t *testing.T, router *mux.Router, refreshToken string) (int, RefreshTokenResponse) {
		req, err := http.NewRequest(http.MethodPost, common.AuthV1+"/refresh", nil)
		require.NoError(t, err)
		req.Header.Set(AuthenticationHeader, refreshToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		defer rr.Result().Body.Close()
		var resp RefreshTokenResponse
		err = json.NewDecoder(rr.Result().Body).Decode(&resp)
		require.NoError(t, err)

		return rr.Result().StatusCode, resp
	}

	// Set the token expiry to something short, so we can get an "expired" token
	accessTokenDuration = time.Duration(1 * time.Second)

	// Register a new user
	t_registerUser(t, router, RegisterNewUserRequest{Username: username, Password: password})

	// Successfully login
	status, loginResp := t_login(t, router, LoginRequest{Username: username, Password: password})
	require.Equal(t, http.StatusOK, status)

	// Hit a protected route to prove our token is valid
	validProtectedRoute(t, router, loginResp.AccessToken)

	// Wait a bit for the token to expire (not great I know)
	time.Sleep(2 * time.Second)

	// Token should be invalid by now, but our refresh token shouldnt
	unauthorizedProtectedRoute(t, router, loginResp.AccessToken)

	// Refresh the token
	status, refreshResponse := t_refresh(t, router, loginResp.RefreshToken)
	require.Equal(t, http.StatusOK, status)

	// Should be able to use the new token to hit the protected route
	validProtectedRoute(t, router, refreshResponse.AccessToken)
}
