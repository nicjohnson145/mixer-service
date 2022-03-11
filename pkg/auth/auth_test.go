package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func newDB(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("auth-testing.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	err = autoMigrate(db)
	require.NoError(t, err)

	cleanup := func() {
		os.Remove("auth-testing.db")
	}

	return db, cleanup
}

func TestRegisterLogin(t *testing.T) {

	body := func() io.Reader {
		return bytes.NewReader([]byte(`{"username": "foo", "password": "bar"}`))
	}

	t.Run("register_login_happy", func(t *testing.T) {
		db, cleanup := newDB(t)
		defer cleanup()

		registerHandler := registerNewUser(db)
		loginHandler := login(db)

		// Register the user
		registerReq, err := http.NewRequest("POST", "/register-user", body())
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		registerHandler(rr, registerReq)
		require.Equal(t, http.StatusOK, rr.Result().StatusCode)

		// Now we should be able to login as that user
		loginReq, err := http.NewRequest("POST", "/login", body())
		require.NoError(t, err)
		rr = httptest.NewRecorder()
		loginHandler(rr, loginReq)
		require.Equal(t, http.StatusOK, rr.Result().StatusCode)
	})

	t.Run("double_register_errors", func(t *testing.T) {
		db, cleanup := newDB(t)
		defer cleanup()

		registerHandler := registerNewUser(db)

		// Register the user
		registerReq, err := http.NewRequest("POST", "/register-user", body())
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		registerHandler(rr, registerReq)
		require.Equal(t, http.StatusOK, rr.Result().StatusCode)

		// Try and register the same user again should result in an error
		registerReq, err = http.NewRequest("POST", "/register-user", body())
		require.NoError(t, err)
		rr = httptest.NewRecorder()
		registerHandler(rr, registerReq)

		require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)

		var response RegisterNewUserResponse
		err = json.NewDecoder(rr.Result().Body).Decode(&response)
		require.NoError(t, err)
		expected := RegisterNewUserResponse{
			Error:   "user foo already registered",
			Success: false,
		}
		require.Equal(t, expected, response)
	})

	t.Run("protected_no_token_rejected", func(t *testing.T) {
		const return_val = "protected endpoint"

		protected := func(w http.ResponseWriter, r *http.Request, c Claims) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, return_val)
		}

		protectedHandler := Protected(protected)

		protectedRequest, err := http.NewRequest("GET", "/protected", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		protectedHandler(rr, protectedRequest)

		require.Equal(t, http.StatusUnauthorized, rr.Result().StatusCode)
	})

	t.Run("protected_valid_token_accepted", func(t *testing.T) {
		const return_val = "protected endpoint"

		protected := func(w http.ResponseWriter, r *http.Request, c Claims) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, return_val)
		}

		protectedHandler := Protected(protected)

		protectedRequest, err := http.NewRequest("GET", "/protected", nil)
		token, err := generateTokenString(TokenInputs{Username: "foobar"})
		require.NoError(t, err)
		protectedRequest.Header.Set(AuthenticationHeader, token)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		protectedHandler(rr, protectedRequest)

		require.Equal(t, http.StatusOK, rr.Result().StatusCode)
		got, err := ioutil.ReadAll(rr.Result().Body)
		require.NoError(t, err)
		require.Equal(t, return_val, string(got))
	})
}
