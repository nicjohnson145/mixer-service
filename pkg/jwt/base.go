package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	AuthenticationHeader = "MixerAuth"
	tokenTypeRefresh     = "refresh-token"
	tokenTypeAccess      = "access-token"
)

var jwtSecret = getSecretKey()
var accessTokenDuration = getAccessTokenDuration()
var refreshTokenDuration = getRefreshTokenDuration()

func SetAccessTokenDuration(t time.Duration) {
	accessTokenDuration = t
}

func GetAccessTokenDuration() time.Duration {
	return accessTokenDuration
}

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
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
		log.Infof("%v set to %v", key, d)
		return d
	} else {
		log.Infof("Using default %v of %v", key, defaultDuration)
		return defaultDuration
	}
}

func GenerateAccessToken(i TokenInputs) (string, error) {
	return generateTokenStringWithExpiry(i, tokenTypeAccess, accessTokenDuration)
}

func GenerateRefreshToken(i TokenInputs) (string, error) {
	return generateTokenStringWithExpiry(i, tokenTypeRefresh, refreshTokenDuration)
}

func generateTokenStringWithExpiry(i TokenInputs, tokenType string, expiry time.Duration) (string, error) {
	claims := &Claims{
		Username:  i.Username,
		TokenType: tokenType,
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

func ValidateRefreshToken(token string) (Claims, error) {
	claims, err := validateToken(token)
	if err != nil {
		return Claims{}, err
	}

	if claims.TokenType != tokenTypeRefresh {
		return Claims{}, fmt.Errorf("token is not refresh token")
	}

	return claims, nil
}

func ValidateAccessToken(token string) (Claims, error) {
	claims, err := validateToken(token)
	if err != nil {
		return Claims{}, err
	}

	if claims.TokenType != tokenTypeAccess {
		return Claims{}, fmt.Errorf("token is not access token")
	}

	return claims, nil
}
