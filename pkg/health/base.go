package health

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"net/http"
	"fmt"
)

func Init(r *mux.Router, db *sql.DB) error {
	defineRoutes(r, db)
	return nil
}

func defineRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc(common.HealthV1, healthCheck)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}
