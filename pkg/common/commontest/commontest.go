package commontest

import (
	"database/sql"
	"net/http"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/db"
	"github.com/stretchr/testify/require"
)

func newDB(t *testing.T, name string) (*sql.DB, func()) {
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

func SetupDbAndRouter(t *testing.T, name string, routeFunc func(*fiber.App, *sql.DB)) (*fiber.App, func()) {
	db, cleanup := newDB(t, name)
	app := common.NewApp()
	routeFunc(app, db)
	return app, cleanup
}

func SetJsonHeader(r *http.Request) {
	r.Header.Set("Content-type", "application/json")
}

func RequireOkStatus(t *testing.T, resp *http.Response) {
	t.Helper()
	require.True(t, resp.StatusCode >= 200 && resp.StatusCode <= 299)
}

func RequireNotOkStatus(t *testing.T, resp *http.Response) {
	t.Helper()
	require.False(t, resp.StatusCode >= 200 && resp.StatusCode <= 299)
}
