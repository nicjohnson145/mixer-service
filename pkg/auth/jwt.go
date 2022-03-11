package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"net/http"
	"errors"
	"fmt"
	"encoding/json"
)

const (
	AuthenticationHeader = "MixerAuth"
)

var jwtSecret = []byte("my-super-secert-key")

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type TokenInputs struct {
	Username string
}

func generateTokenString(i TokenInputs) (string, error) {
	claims := &Claims{
		Username: i.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateToken(token string) (Claims, error) {
	claims := Claims{}

	tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return Claims{}, err
	}

	if !tkn.Valid {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}


func Protected(handler ClaimsHttpHandler) HttpHandler {

	writeUnauthorized := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusUnauthorized)
		bytes, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		fmt.Fprintln(w, string(bytes))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get(AuthenticationHeader)
		if val == "" {
			writeUnauthorized(w)
			return
		}

		claims, err := validateToken(val)
		if err != nil {
			writeUnauthorized(w)
			return
		}

		handler(w, r, claims)
	}
}
