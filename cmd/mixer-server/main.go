package main

import (
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/db"
	"github.com/nicjohnson145/mixer-service/pkg/drink"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	db := db.NewDBOrDie(common.DefaultedEnvVar("DB_PATH", "mixer.db"))

	if err := auth.Init(r, db); err != nil {
		log.Fatal(err)
	}

	if err := drink.Init(r, db); err != nil {
		log.Fatal(err)
	}

	port := common.DefaultedEnvVar("PORT", "30000")

	log.Info("Listening on port ", port)
	http.ListenAndServe(":"+port, r)
}
