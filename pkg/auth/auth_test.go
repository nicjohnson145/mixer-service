package auth

import (
	"testing"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/sqlite"
	"github.com/stretchr/testify/require"
	"os"
	"net/http"
	"net/http/httptest"
	"bytes"
	"io"
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

	t.Run("register_login_happy", func(t *testing.T) {
		db, cleanup := newDB(t)
		defer cleanup()

		registerHandler := registerNewUser(db)
		loginHandler := login(db)

		body := func() io.Reader {
			return bytes.NewReader([]byte(`{"username": "foo", "password": "bar"}`))
		}

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
}
