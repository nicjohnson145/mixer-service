package auth

import (
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type User struct {
	Username string `gorm:"primaryKey"`
	Password string
}

type httpHandler func(http.ResponseWriter, *http.Request)

type basicResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

func Init(r *mux.Router, db *gorm.DB) error {
	if err := autoMigrate(db); err != nil {
		return err
	}

	defineRoutes(r, db)
	return nil
}

func autoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}

	return nil
}

func defineRoutes(r *mux.Router, db *gorm.DB) {
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
