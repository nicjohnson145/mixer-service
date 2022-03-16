package auth

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RefreshTokenResponse struct {
	Error       string `json:"error,omitempty"`
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token,omitempty"`
}

func refresh() ClaimsHttpHandler {

	writeResponse := func(w http.ResponseWriter, status int, error string, token string) {
		w.WriteHeader(status)
		bytes, _ := json.Marshal(RefreshTokenResponse{
			Error:       error,
			Success:     status >= 200 && status <= 299,
			AccessToken: token,
		})
		fmt.Fprintln(w, string(bytes))
	}

	writeInternalError := func(w http.ResponseWriter, err error, location string) {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"operation": location,
		}).Error("error during user login")
		writeResponse(w, http.StatusInternalServerError, "internal error", "")
	}

	writeSucess := func(w http.ResponseWriter, accessToken string) {
		writeResponse(w, http.StatusOK, "", accessToken)
	}

	return func(w http.ResponseWriter, r *http.Request, claims Claims) {
		newToken, err := GenerateAccessToken(TokenInputs{Username: claims.Username})
		if err != nil {
			writeInternalError(w, err, "generating new access token")
			return
		}

		writeSucess(w, newToken)
	}
}
