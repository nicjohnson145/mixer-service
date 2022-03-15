package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

const (
	AuthenticationHeader = "MixerAuth"
)

var jwtSecret = getSecretKey()

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type TokenInputs struct {
	Username string
}

func getSecretKey() []byte {
	if val, ok := os.LookupEnv("JWT_SECRET"); ok {
		return []byte(val)
	} else {
		log.Warning("No JWT_SECRET set, defaulting to hardcoded secret. THIS IS INSECURE!!")
		return []byte("super-secret-jwt-key")
	}
}

func GenerateTokenString(i TokenInputs) (string, error) {
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

func RequiresValidToken(handler ClaimsHttpHandler) common.HttpHandler {

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
