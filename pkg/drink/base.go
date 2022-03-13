package drink

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"net/http"
	"strings"
)

type Drink struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	Username       string   `json:"username"`
	PrimaryAlcohol string   `json:"primary_alcohol"`
	PreferredGlass string   `json:"preferred_glass"`
	Ingredients    []string `json:"ingredients"`
	Instructions   string   `json:"instructions"`
	Notes          string   `json:"notes"`
}

func Init(r *mux.Router, db *sql.DB) error {
	defineRoutes(r, db)
	return nil
}

func defineRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc(common.DrinksV1+"/create", auth.Protected(createDrink(db))).Methods(http.MethodPost)
	r.HandleFunc(common.DrinksV1+"/{id:[0-9]+}", auth.Protected(getDrink(db))).Methods(http.MethodGet)
}

func fromDb(d Model) (Drink, error) {
	ingredients, err := fromCSV(d.Ingredients)
	if err != nil {
		return Drink{}, err
	}

	return Drink{
		Name:           d.Name,
		Username:       d.Username,
		PrimaryAlcohol: d.PrimaryAlcohol,
		PreferredGlass: d.PreferredGlass,
		Ingredients:    ingredients,
		Instructions:   d.Instructions,
		Notes:          d.Notes,
	}, nil
}

//func toDb(d Drink) (Model, error) {
//    ingredients, err := toCSV(d.Ingredients)
//    if err != nil {
//        return Model{}, err
//    }

//    return Model{
//        Name:           d.Name,
//        Username:       d.Username,
//        PrimaryAlcohol: d.PrimaryAlcohol,
//        PreferredGlass: d.PreferredGlass,
//        Ingredients:    ingredients,
//        Instructions:   d.Instructions,
//        Notes:          d.Notes,
//    }, nil
//}

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
