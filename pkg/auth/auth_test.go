package auth

import (
	"database/sql"
	"encoding/json"
	"github.com/nicjohnson145/mixer-service/pkg/db"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"github.com/gorilla/mux"
)

func newDB(t *testing.T) (*sql.DB, func()) {
	const name = "auth.db"
	db, err := db.NewDB(name)
	require.NoError(t, err)

	cleanup := func() {
		err := os.Remove(name)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
	return db, cleanup
}

func TestRegisterLogin(t *testing.T) {
	// Suppress log output by default
	log.SetOutput(ioutil.Discard)

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

	db, cleanup := newDB(t)
	defer cleanup()

	router := mux.NewRouter()
	defineRoutes(router, db)

	realUser := func() io.Reader {
		return strings.NewReader(`{"username": "foo", "password": "bar"}`)
	}

	// Register the user
	registerReq, err := http.NewRequest("POST", common.AuthV1 + "/register-user", realUser())
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, registerReq)
	require.Equal(t, http.StatusOK, rr.Result().StatusCode)

	for _, tc := range loginData {
		t.Run("login_cases_"+tc.name, func(t *testing.T) {
			loginReq, err := http.NewRequest("POST", common.AuthV1 + "/login", strings.NewReader(tc.input))
			require.NoError(t, err)
			rr = httptest.NewRecorder()
			router.ServeHTTP(rr, loginReq)
			require.Equal(t, tc.expectedCode, rr.Result().StatusCode)

			defer rr.Result().Body.Close()
			var resp LoginResponse
			err = json.NewDecoder(rr.Result().Body).Decode(&resp)
			require.NoError(t, err)
			if tc.expectedToken {
				require.NotEmpty(t, resp.Token)
			} else {
				require.Empty(t, resp.Token)
			}
		})
	}
}
