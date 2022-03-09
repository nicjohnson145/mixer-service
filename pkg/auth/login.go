package auth

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"errors"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	basicResponse
}


func login(db *gorm.DB) httpHandler {

	writeLoginResponse := func (w http.ResponseWriter, status int, error string) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(LoginResponse{
			basicResponse: basicResponse{
				Error: error,
				Success: status >= 200 && status <= 299,
			},
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeUnauthorizedError := func (w http.ResponseWriter) {
		writeLoginResponse(w, http.StatusUnauthorized, "unauthorized")
	}

	writeBadRequestError := func(w http.ResponseWriter, msg string) {
		writeLoginResponse(w, http.StatusBadRequest, msg)
	}

	writeInternalError := func(w http.ResponseWriter) {
		writeLoginResponse(w, http.StatusInternalServerError, "internal error")
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
		
		w.WriteHeader(http.StatusOK)
	}
}
