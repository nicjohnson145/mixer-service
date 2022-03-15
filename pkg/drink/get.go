package drink

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type GetDrinkResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
	Drink   *Drink `json:"drink"`
}

func getDrink(db *sql.DB) auth.ClaimsHttpHandler {
	writeResponse := func(w http.ResponseWriter, msg string, status int, d *Drink) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(GetDrinkResponse{
			Error:   msg,
			Success: status >= 200 && status <= 299,
			Drink:   d,
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeInternalError := func(w http.ResponseWriter, err error, location string, id int64) {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"operation": location,
			"id":        id,
		}).Error("internal error while getting drink")
		writeResponse(w, err.Error(), http.StatusInternalServerError, nil)
	}

	writeNotFound := func(w http.ResponseWriter) {
		writeResponse(w, "not found", http.StatusNotFound, nil)
	}

	writeSucess := func(w http.ResponseWriter, d Drink) {
		writeResponse(w, "", http.StatusOK, &d)
	}

	return func(w http.ResponseWriter, r *http.Request, claims auth.Claims) {
		vars := mux.Vars(r)
		// Mux handles that the route is a number
		id, _ := strconv.ParseInt(vars["id"], 10, 64)

		model, err := getByID(id, db)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				writeNotFound(w)
				return
			} else {
				writeInternalError(w, err, "getting drink", id)
				return
			}
		}

		drink, err := fromDb(*model)
		if err != nil {
			writeInternalError(w, err, "converting from DB type", id)
			return
		}

		if drink.Username != claims.Username && drink.Publicity != DrinkPublicityPublic {
			writeNotFound(w)
			return
		}

		writeSucess(w, drink)
	}
}
