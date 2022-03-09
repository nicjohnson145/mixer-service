package main

import (
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"
)

func main() {
	r := mux.NewRouter()
	db, err := gorm.Open(sqlite.Open("foo-bar.db"), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := auth.Init(r, db); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	port := "30000"
	log.Info("Listening on port ", port)
	http.ListenAndServe(":" + port, r)
}
