package auth

import (
	"database/sql"
	"encoding/json"
	"github.com/nicjohnson145/mixer-service/pkg/db"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func newDB(t *testing.T) (*sql.DB, func()) {
	db, err := db.NewDB("auth.db")
	require.NoError(t, err)

	cleanup := func() {
		err := os.Remove("auth.db")
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

	registerHandler := registerNewUser(db)
	loginHandler := login(db)

	realUser := func() io.Reader {
		return strings.NewReader(`{"username": "foo", "password": "bar"}`)
	}

	// Register the user
	registerReq, err := http.NewRequest("POST", "/register-user", realUser())
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	registerHandler(rr, registerReq)
	require.Equal(t, http.StatusOK, rr.Result().StatusCode)

	for _, tc := range loginData {
		t.Run("login_cases_"+tc.name, func(t *testing.T) {
			loginReq, err := http.NewRequest("POST", "/login", strings.NewReader(tc.input))
			require.NoError(t, err)
			rr = httptest.NewRecorder()
			loginHandler(rr, loginReq)
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
