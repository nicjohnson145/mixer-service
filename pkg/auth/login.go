package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

func login(db *gorm.DB) httpHandler {

	writeLoginResponse := func(w http.ResponseWriter, status int, error string) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(LoginResponse{
			Error:   error,
			Success: status >= 200 && status <= 299,
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeUnauthorizedError := func(w http.ResponseWriter) {
		writeLoginResponse(w, http.StatusUnauthorized, "unauthorized")
	}

	writeBadRequestError := func(w http.ResponseWriter, msg string) {
		writeLoginResponse(w, http.StatusBadRequest, msg)
	}

	writeInternalError := func(w http.ResponseWriter) {
		writeLoginResponse(w, http.StatusInternalServerError, "internal error")
	}

	writeSucess := func(w http.ResponseWriter) {
		writeLoginResponse(w, http.StatusOK, "")
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
				writeInternalError(w)
				return
			}
		}

		if !comparePasswords(existingUser.Password, payload.Password) {
			writeUnauthorizedError(w)
			return
		}

		writeSucess(w)
	}
}
