package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	log "github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
}

func login(db *gorm.DB) httpHandler {

	writeLoginResponse := func(w http.ResponseWriter, status int, error string, token string) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(LoginResponse{
			Error:   error,
			Success: status >= 200 && status <= 299,
			Token: token,
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeUnauthorizedError := func(w http.ResponseWriter) {
		writeLoginResponse(w, http.StatusUnauthorized, "unauthorized", "")
	}

	writeBadRequestError := func(w http.ResponseWriter, msg string) {
		writeLoginResponse(w, http.StatusBadRequest, msg, "")
	}

	writeInternalError := func(w http.ResponseWriter, err error) {
		log.WithField("error", err.Error()).Error("Internal server error")
		writeLoginResponse(w, http.StatusInternalServerError, "internal error", "")
	}

	writeSucess := func(w http.ResponseWriter, token string) {
		writeLoginResponse(w, http.StatusOK, "", token)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload LoginRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		defer r.Body.Close()

		if err != nil {
			writeBadRequestError(w, err.Error())
			return
		}
		var existingUser User
		result := db.Model(&User{}).First(&existingUser, "username = ?", payload.Username)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				writeUnauthorizedError(w)
				return
			} else {
				writeInternalError(w, err)
				return
			}
		}

		if !comparePasswords(existingUser.Password, payload.Password) {
			writeUnauthorizedError(w)
			return
		}

		tokenStr, err := generateTokenString(TokenInputs{Username: payload.Username})
		if err != nil {
			writeInternalError(w, err)
			return
		}

		writeSucess(w, tokenStr)
	}
}
