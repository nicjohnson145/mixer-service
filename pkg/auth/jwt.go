package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
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
var accessTokenDuration = getAccessTokenDuration()
var refreshTokenDuration = getRefreshTokenDuration()

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

func getAccessTokenDuration() time.Duration {
	return lookupDefaultedDuration("ACCESS_TOKEN_DURATION", time.Duration(5*time.Minute))
}

func getRefreshTokenDuration() time.Duration {
	// Defaults to ~1 month
	return lookupDefaultedDuration("REFRESH_TOKEN_DURATION", time.Duration(730*time.Hour))
}

func lookupDefaultedDuration(key string, defaultDuration time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		d, err := time.ParseDuration(val)
		if err != nil {
			log.Fatal(fmt.Sprintf("error parsing %v: %v", key, err))
		}
		return d
	} else {
		log.Info(fmt.Sprintf("Using default %v of %v", key, defaultDuration))
		return defaultDuration
	}
}

func GenerateAccessToken(i TokenInputs) (string, error) {
	return GenerateTokenStringWithExpiry(i, accessTokenDuration)
}

func GenerateRefreshToken(i TokenInputs) (string, error) {
	return GenerateTokenStringWithExpiry(i, refreshTokenDuration)
}

func GenerateTokenStringWithExpiry(i TokenInputs, expiry time.Duration) (string, error) {
	claims := &Claims{
		Username: i.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiry).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateToken(token string) (Claims, error) {
	claims := Claims{}

	tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
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
