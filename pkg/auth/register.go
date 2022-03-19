package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RegisterNewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterNewUserResponse struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

func registerNewUser(db *sql.DB) ClaimsHttpHandler {

	writeRegisterNewUserReponse := func(w http.ResponseWriter, status int, msg string) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(RegisterNewUserResponse{
			Error:   msg,
			Success: status >= 200 && status <= 299,
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeInternalError := func(w http.ResponseWriter, err error, operation string) {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"operation": operation,
		}).Error("error registering new user")
		writeRegisterNewUserReponse(w, http.StatusInternalServerError, "internal error")
	}

	writeBadRequestError := func(w http.ResponseWriter, msg string) {
		writeRegisterNewUserReponse(w, http.StatusBadRequest, msg)
	}

	writeSucess := func(w http.ResponseWriter) {
		writeRegisterNewUserReponse(w, http.StatusOK, "")
	}

	return func(w http.ResponseWriter, r *http.Request, claims Claims) {

		var payload RegisterNewUserRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		defer r.Body.Close()

		if err != nil {
			writeBadRequestError(w, err.Error())
			return
		}

		existingUser, err := getUserByName(payload.Username, db)
		if err != nil && !errors.Is(err, common.ErrNotFound) {
			writeInternalError(w, err, "checking for existing user")
			return
		}

		if existingUser != nil {
			writeBadRequestError(w, fmt.Sprintf("user %v already exists", payload.Username))
			return
		}

		hashedPw, err := hashPassword(payload.Password)
		if err != nil {
			writeInternalError(w, err, "hashing password")
			return
		}

		newUser := UserModel{
			Username: payload.Username,
			Password: hashedPw,
		}
		err = createUser(newUser, db)
		if err != nil {
			writeInternalError(w, err, "inserting into db")
			return
		}

		writeSucess(w)
	}

}
