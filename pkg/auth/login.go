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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Error        string `json:"error,omitempty"`
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func login(db *sql.DB) common.HttpHandler {

	writeLoginResponse := func(w http.ResponseWriter, status int, error string, token string, refreshToken string) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(LoginResponse{
			Error:        error,
			Success:      status >= 200 && status <= 299,
			AccessToken:  token,
			RefreshToken: refreshToken,
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeUnauthorizedError := func(w http.ResponseWriter, user string, reason string) {
		log.WithFields(log.Fields{
			"user":   user,
			"reason": reason,
		}).Info("invalid login attempt")
		writeLoginResponse(w, http.StatusUnauthorized, "unauthorized", "", "")
	}

	writeBadRequestError := func(w http.ResponseWriter, msg string) {
		writeLoginResponse(w, http.StatusBadRequest, msg, "", "")
	}

	writeInternalError := func(w http.ResponseWriter, err error, location string) {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"operation": location,
		}).Error("error during user login")
		writeLoginResponse(w, http.StatusInternalServerError, "internal error", "", "")
	}

	writeSucess := func(w http.ResponseWriter, accessToken string, refreshToken string) {
		writeLoginResponse(w, http.StatusOK, "", accessToken, refreshToken)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload LoginRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		defer r.Body.Close()

		if err != nil {
			writeBadRequestError(w, err.Error())
			return
		}
		existingUser, err := getUserByName(payload.Username, db)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				writeUnauthorizedError(w, payload.Username, "fetching from db")
				return
			} else {
				writeInternalError(w, err, "getting user from db")
				return
			}
		}

		if !comparePasswords(existingUser.Password, payload.Password) {
			writeUnauthorizedError(w, payload.Username, "comparing passwords")
			return
		}

		accessStr, err := GenerateAccessToken(TokenInputs{Username: payload.Username})
		if err != nil {
			writeInternalError(w, err, "generating access token")
			return
		}

		refreshStr, err := GenerateAccessToken(TokenInputs{Username: payload.Username})
		if err != nil {
			writeInternalError(w, err, "generating refresh token")
			return
		}

		writeSucess(w, accessStr, refreshStr)
	}
}
