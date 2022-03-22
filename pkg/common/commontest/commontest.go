package commontest

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/db"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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
	app := fiber.New()
	routeFunc(app, db)
	return app, cleanup
}
