package auth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupDbAndRouter(t *testing.T) (*mux.Router, func()) {
	return commontest.SetupDbAndRouter(t, "auth.db", defineRoutes)
}

func TestRegisterLogin(t *testing.T) {
	loginData := []struct {
		name          string
		input         string
		expectedCode  int
		expectedToken bool
	}{
		{
			name:          "missing_user",
			input:         `{"username": "not_a_user", "password": "bar"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedToken: false,
		},
		{
			name:          "incorrect_password",
			input:         `{"username": "foo", "password": "wrong_password"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedToken: false,
		},
		{
			name:          "valid_login",
			input:         `{"username": "foo", "password": "bar"}`,
			expectedCode:  http.StatusOK,
			expectedToken: true,
		},
	}

	router, cleanup := setupDbAndRouter(t)
	defer cleanup()

	realUser := func() io.Reader {
		return strings.NewReader(`{"username": "foo", "password": "bar"}`)
	}

	// Register the user
	registerReq, err := http.NewRequest("POST", common.AuthV1+"/register-user", realUser())
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, registerReq)
	require.Equal(t, http.StatusOK, rr.Result().StatusCode)

	for _, tc := range loginData {
		t.Run("login_cases_"+tc.name, func(t *testing.T) {
			loginReq, err := http.NewRequest("POST", common.AuthV1+"/login", strings.NewReader(tc.input))
			require.NoError(t, err)
			rr = httptest.NewRecorder()
			router.ServeHTTP(rr, loginReq)
			require.Equal(t, tc.expectedCode, rr.Result().StatusCode)

			defer rr.Result().Body.Close()
			var resp LoginResponse
			err = json.NewDecoder(rr.Result().Body).Decode(&resp)
			require.NoError(t, err)
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
