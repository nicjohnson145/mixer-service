package auth

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"errors"
)

type RegisterNewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterNewUserResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

func registerNewUser(db *gorm.DB) httpHandler {

	writeRegisterNewUserReponse := func (w http.ResponseWriter, status int, msg string) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(RegisterNewUserResponse{
			Error: msg,
			Success: status >= 200 && status <= 299,
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeInternalError := func(w http.ResponseWriter) {
		writeRegisterNewUserReponse(w, http.StatusInternalServerError, "internal error")
	}

	writeBadRequestError := func(w http.ResponseWriter, msg string) {
		writeRegisterNewUserReponse(w, http.StatusBadRequest, msg)
	}

	writeSucess := func(w http.ResponseWriter) {
		writeRegisterNewUserReponse(w, http.StatusOK, "")
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var payload RegisterNewUserRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		defer r.Body.Close()

		if err != nil {
			writeBadRequestError(w, err.Error())
			return
		}

		var existingUser User
		result := db.Model(&User{}).First(&existingUser, "username = ?", payload.Username)

		// If we found something, return bad request, user exists
		if result.Error == nil {
			writeBadRequestError(w, fmt.Sprintf("user %v already registered", payload.Username))
			return
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			writeInternalError(w)
			return
		}

		hashedPw, err := hashPassword(payload.Password)
		if err != nil {
			writeInternalError(w)
			return
		}

		result = db.Create(&User{Username: payload.Username, Password: hashedPw})
		if result.Error != nil {
			writeInternalError(w)
			return
		}

		writeSucess(w)
	}

}

