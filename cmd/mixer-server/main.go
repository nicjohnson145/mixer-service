package main

import (
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/db"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	db := db.NewDBOrDie("foo-bar.db")

	if err := auth.Init(r, db); err != nil {
		log.Fatal(err)
	}

	port := "30000"
	log.Info("Listening on port ", port)
	http.ListenAndServe(":"+port, r)
}
