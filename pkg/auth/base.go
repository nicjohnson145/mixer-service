package auth

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type User struct {
	Username string
	Password string
}

type ClaimsHttpHandler func(http.ResponseWriter, *http.Request, Claims)

func Init(r *mux.Router, db *sql.DB) error {
	defineRoutes(r, db)
	return nil
}

func defineRoutes(r *mux.Router, db *sql.DB) {
	if common.DefaultedEnvVar("PROTECT_REGISTER_ENDPOINT", "false") == "true" {
		r.HandleFunc(common.AuthV1+"/register-user", RequiresValidAccessToken(registerNewUser(db))).Methods(http.MethodPost)
	} else {
		r.HandleFunc(common.AuthV1+"/register-user", noopProtection(registerNewUser(db))).Methods(http.MethodPost)
	}
	r.HandleFunc(common.AuthV1+"/login", login(db)).Methods(http.MethodPost)
	r.HandleFunc(common.AuthV1+"/refresh", requiresValidRefreshToken(refresh())).Methods(http.MethodPost)
}

func hashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func comparePasswords(hashedPw string, plainPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(plainPw))
	return err == nil
}

func noopProtection(handler ClaimsHttpHandler) common.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, Claims{})
	}
}
