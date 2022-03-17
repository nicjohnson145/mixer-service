package drink

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"net/http"
)

const (
	DrinkPublicityPublic  = "public"
	DrinkPublicityPrivate = "private"
)

type Drink struct {
	ID       int64  `json:"id"`
	Username string `json:"username" validate:"required"`
	drinkData
}

var validate = validator.New()

func Init(r *mux.Router, db *sql.DB) error {
	defineRoutes(r, db)
	return nil
}

func defineRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc(common.DrinksV1+"/create", auth.RequiresValidAccessToken(createDrink(db))).Methods(http.MethodPost)
	r.HandleFunc(common.DrinksV1+"/{id:[0-9]+}", auth.RequiresValidAccessToken(getDrink(db))).Methods(http.MethodGet)
	r.HandleFunc(common.DrinksV1+"/{id:[0-9]+}", auth.RequiresValidAccessToken(deleteDrink(db))).Methods(http.MethodDelete)
	r.HandleFunc(common.DrinksV1+"/{id:[0-9]+}", auth.RequiresValidAccessToken(updateDrink(db))).Methods(http.MethodPut)
	r.HandleFunc(common.DrinksV1+"/by-user/{username}", auth.RequiresValidAccessToken(getDrinksByUser(db))).Methods(http.MethodGet)
}

