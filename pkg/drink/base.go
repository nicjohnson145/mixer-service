package drink

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"net/http"
	"strings"
)

const (
	DrinkPublicityPublic  = "public"
	DrinkPublicityPrivate = "private"
)

type drinkData struct {
	Name           string   `json:"name" validate:"required"`
	PrimaryAlcohol string   `json:"primary_alcohol" validate:"required"`
	PreferredGlass string   `json:"preferred_glass"`
	Ingredients    []string `json:"ingredients" validate:"required"`
	Instructions   string   `json:"instructions"`
	Notes          string   `json:"notes"`
	Publicity      string   `json:"publicity" validate:"required"`
}

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
	r.HandleFunc(common.DrinksV1+"/create", auth.RequiresValidToken(createDrink(db))).Methods(http.MethodPost)
	r.HandleFunc(common.DrinksV1+"/{id:[0-9]+}", auth.RequiresValidToken(getDrink(db))).Methods(http.MethodGet)
	r.HandleFunc(common.DrinksV1+"/{id:[0-9]+}", auth.RequiresValidToken(deleteDrink(db))).Methods(http.MethodDelete)
	r.HandleFunc(common.DrinksV1+"/{id:[0-9]+}", auth.RequiresValidToken(updateDrink(db))).Methods(http.MethodPut)
	r.HandleFunc(common.DrinksV1+"/by-user/{username}", auth.RequiresValidToken(getDrinksByUser(db))).Methods(http.MethodGet)
}

func fromDb(d Model) (Drink, error) {
	ingredients, err := fromCSV(d.Ingredients)
	if err != nil {
		return Drink{}, err
	}

	return Drink{
		ID:       d.ID,
		Username: d.Username,
		drinkData: drinkData{
			Name:           d.Name,
			PrimaryAlcohol: d.PrimaryAlcohol,
			PreferredGlass: d.PreferredGlass,
			Ingredients:    ingredients,
			Instructions:   d.Instructions,
			Notes:          d.Notes,
			Publicity:      d.Publicity,
		},
	}, nil
}

func toDb(d Drink) (Model, error) {
	ingredients, err := toCSV(d.Ingredients)
	if err != nil {
		return Model{}, err
	}

	return Model{
		ID:             d.ID,
		Name:           d.Name,
		Username:       d.Username,
		PrimaryAlcohol: d.PrimaryAlcohol,
		PreferredGlass: d.PreferredGlass,
		Ingredients:    ingredients,
		Instructions:   d.Instructions,
		Notes:          d.Notes,
		Publicity:      d.Publicity,
	}, nil
}

func toCSV(s []string) (string, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	err := w.Write(s)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.ReplaceAll(buf.String(), "\n", ""), nil
}

func fromCSV(s string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	return r.Read()
}
