package commontest

import (
	"database/sql"
	"github.com/gorilla/mux"
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

func SetupDbAndRouter(t *testing.T, name string, routeFunc func(*mux.Router, *sql.DB)) (*mux.Router, func()) {
	db, cleanup := newDB(t, name)
	router := mux.NewRouter()
	routeFunc(router, db)
	return router, cleanup
}
