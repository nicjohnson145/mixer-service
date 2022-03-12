package drink

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type CreateDrinkRequest struct {
	Name           string   `json:"name"`
	PrimaryAlcohol string   `json:"primary_alcohol"`
	PreferredGlass string   `json:"preferred_glass"`
	Ingredients    []string `json:"ingredients"`
	Instructions   string   `json:"instructions"`
	Notes          string   `json:"notes"`
}

type CreateDrinkResponse struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
	ID      int64  `json:"id,omitempty"`
}

func createDrink(db *sql.DB) auth.ClaimsHttpHandler {

	writeResponse := func(w http.ResponseWriter, msg string, status int, id int64) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(CreateDrinkResponse{
			Error:   msg,
			Success: status >= 200 && status <= 299,
			ID: id,
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeBadRequest := func(w http.ResponseWriter, msg string, location string) {
		log.WithFields(log.Fields{
			"error": msg,
			"operation": location,
		}).Error("bad create drink request")
		writeResponse(w, msg, http.StatusBadRequest, 0)
	}

	writeInternalError := func(w http.ResponseWriter, err error, location string) {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"operation": location,
		}).Error("internal error during drink creation")
		writeResponse(w, err.Error(), http.StatusInternalServerError, 0)
	}

	writeSucess := func(w http.ResponseWriter, id int64) {
		writeResponse(w, "success", http.StatusOK, id)
	}

	return func(w http.ResponseWriter, r *http.Request, claims auth.Claims) {
		var payload CreateDrinkRequest
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			writeBadRequest(w, err.Error(), "decoding payload")
			return
		}

		existingDrink, err := getByNameAndUsername(payload.Name, claims.Username, db)
		if err != nil && !errors.Is(err, common.ErrNotFound) {
			writeInternalError(w, err, "getting existing drink")
			return
		}

		if existingDrink != nil {
			writeBadRequest(
				w,
				fmt.Sprintf("user %v already has a drink named %v", claims.Username, payload.Name),
				"getting existing drink",
			)
			return
		}

		ingredients, err := toCSV(payload.Ingredients)
		if err != nil {
			writeBadRequest(w, err.Error(), "converting to DB model")
			return
		}

		model := Model{
			Name: payload.Name,
			Username: claims.Username,
			PrimaryAlcohol: payload.PrimaryAlcohol,
			PreferredGlass: payload.PreferredGlass,
			Ingredients: ingredients,
			Instructions: payload.Instructions,
			Notes: payload.Notes,
		}

		id, err := create(model, db)
		if err != nil {
			writeInternalError(w, err, "inserting into db")
			return
		}

		writeSucess(w, id)
	}
}
