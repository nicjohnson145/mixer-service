package drink

import (
	"bytes"
	"encoding/csv"
	"strings"
)

type Drink struct {
	Name           string   `json:"name"`
	Username       string   `json:"username"`
	PrimaryAlcohol string   `json:"primary_alcohol"`
	PreferredGlass string   `json:"preferred_glass"`
	Ingredients    []string `json:"ingredients"`
	Instructions   string   `json:"instructions"`
	Notes          string   `json:"notes"`
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

func toDb(d Drink) (Model, error) {
	ingredients, err := toCSV(d.Ingredients)
	if err != nil {
		return Model{}, err
	}

	return Model{
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
