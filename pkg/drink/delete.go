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

type DeleteDrinkResponse struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

func deleteDrink(db *sql.DB) auth.ClaimsHttpHandler {
	writeResponse := func(w http.ResponseWriter, status int, err string) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(DeleteDrinkResponse{
			Error:   err,
			Success: status >= 200 && status <= 299,
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeNotFound := func(w http.ResponseWriter) {
		writeResponse(w, http.StatusNotFound, "not found")
	}

	writeInternalError := func(w http.ResponseWriter, err error, id int64, operation string) {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"operation": operation,
			"id":        id,
		}).Error("internal error while deleting drink")
		writeResponse(w, http.StatusInternalServerError, "internal error")
	}

	writeSucces := func(w http.ResponseWriter) {
		writeResponse(w, http.StatusOK, "")
	}

	return func(w http.ResponseWriter, r *http.Request, claims auth.Claims) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 64)

		model, err := getByID(id, db)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				writeNotFound(w)
				return
			} else {
				writeInternalError(w, err, id, "fetching by id")
				return
			}
		}
		if model.Username != claims.Username {
			writeNotFound(w)
			return
		}

		err = deleteModel(id, db)
		if err != nil {
			writeInternalError(w, err, id, "deleting drink")
			return
		}

		writeSucces(w)
	}
}
