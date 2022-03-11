package drink

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"github.com/huandu/go-sqlbuilder"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"strings"
)

var DrinkModelStruct = sqlbuilder.NewStruct(new(DrinkModel))

const (
	DrinkTableName = "drink"
)

type DrinkModel struct {
	ID             int    `db:"id"`
	Name           string `db:"name"`
	Username       string `db:"username"`
	PrimaryAlcohol string `db:"primary_alcohol"`
	PreferredGlass string `db:"preferred_glass"`
	Ingredients    string `db:"ingredients"`
	Instructions   string `db:"instructions"`
	Notes          string `db:"notes"`
}

type Drink struct {
	Name           string   `json:"name"`
	Username       string   `json:"username"`
	PrimaryAlcohol string   `json:"primary_alcohol"`
	PreferredGlass string   `json:"preferred_glass"`
	Ingredients    []string `json:"ingredients"`
	Instructions   string   `json:"instructions"`
	Notes          string   `json:"notes"`
}

func fromDb(d DrinkModel) (Drink, error) {
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

func toDb(d Drink) (DrinkModel, error) {
	ingredients, err := toCSV(d.Ingredients)
	if err != nil {
		return DrinkModel{}, err
	}

	return DrinkModel{
		Name:           d.Name,
		Username:       d.Username,
		PrimaryAlcohol: d.PrimaryAlcohol,
		PreferredGlass: d.PreferredGlass,
		Ingredients:    ingredients,
		Instructions:   d.Instructions,
		Notes:          d.Notes,
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

func getDrinkByID(id int, db *sql.DB) (*DrinkModel, error) {
	sb := DrinkModelStruct.SelectFrom(DrinkTableName)
	sb.Where(sb.Equal("id", id))

	sql, args := sb.Build()
	rows, err := db.Query(sql, args...)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	hasRow := rows.Next()
	if !hasRow {
		return nil, common.ErrNotFound
	}

	var drink DrinkModel
	err = rows.Scan(DrinkModelStruct.Addr(&drink)...)
	if err != nil {
		return nil, err
	}

	return &drink, nil
}

