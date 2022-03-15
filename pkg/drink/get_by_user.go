package drink

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type GetDrinksByUserResponse struct {
	Error   string  `json:"error"`
	Success bool    `json:"success"`
	Drinks  []Drink `json:"drinks"`
}

func getDrinksByUser(db *sql.DB) auth.ClaimsHttpHandler {
	writeResponse := func(w http.ResponseWriter, msg string, status int, drinks []Drink) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(GetDrinksByUserResponse{
			Error:   msg,
			Success: status >= 200 && status <= 299,
			Drinks:  drinks,
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeInternalError := func(w http.ResponseWriter, err error, location string, username string) {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"operation": location,
			"username":  username,
		}).Error("internal error while getting drink")
		writeResponse(w, err.Error(), http.StatusInternalServerError, nil)
	}

	writeSucess := func(w http.ResponseWriter, drinks []Drink) {
		writeResponse(w, "", http.StatusOK, drinks)
	}

	return func(w http.ResponseWriter, r *http.Request, claims auth.Claims) {
		vars := mux.Vars(r)
		username := vars["username"]

		var modelList []Model
		var err error
		if username == claims.Username {
			modelList, err = getAllDrinksByUser(username, db)
		} else {
			modelList, err = getAllPublicDrinksByUser(username, db)
		}
		if err != nil {
			writeInternalError(w, err, "getting drinks", username)
			return
		}

		drinks := make([]Drink, 0, len(modelList))
		for _, m := range modelList {
			d, err := fromDb(m)
			if err != nil {
				writeInternalError(w, err, "converting db response", username)
				return
			}
			drinks = append(drinks, d)
		}

		writeSucess(w, drinks)
	}
}
