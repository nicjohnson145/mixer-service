package auth

import (
	"database/sql"
	"github.com/gorilla/mux"
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
	r.HandleFunc("/api/v1/register-user", registerNewUser(db)).Methods("POST")
	r.HandleFunc("/api/v1/login", login(db)).Methods("POST")
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
