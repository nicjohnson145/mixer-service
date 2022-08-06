package drink

import (
	"bytes"
	"encoding/csv"
	"strings"
)

func fromDb(d Model) (Drink, error) {
	ingredients, err := fromCSV(d.Ingredients)
	if err != nil {
		return Drink{}, err
	}

	tags, err := fromCSV(d.Tags)
	if err != nil {
		return Drink{}, err
	}

	return Drink{
		ID:       d.ID,
		Username: d.Username,
		DrinkData: DrinkData{
			Name:             d.Name,
			PrimaryAlcohol:   d.PrimaryAlcohol,
			PreferredGlass:   d.PreferredGlass,
			Ingredients:      ingredients,
			Instructions:     d.Instructions,
			Notes:            d.Notes,
			Publicity:        d.Publicity,
			UnderDevelopment: toBool(d.UnderDevelopment),
			Tags:             tags,
			Favorite:         toBool(d.Favorite),
		},
	}, nil
}

func toDb(d Drink) (Model, error) {
	ingredients, err := toCSV(d.Ingredients)
	if err != nil {
		return Model{}, err
	}

	tags, err := toCSV(d.Tags)
	if err != nil {
		return Model{}, err
	}

	return Model{
		ID:               d.ID,
		Name:             d.Name,
		Username:         d.Username,
		PrimaryAlcohol:   d.PrimaryAlcohol,
		PreferredGlass:   d.PreferredGlass,
		Ingredients:      ingredients,
		Instructions:     d.Instructions,
		Notes:            d.Notes,
		Publicity:        d.Publicity,
		UnderDevelopment: fromBool(d.UnderDevelopment),
		Tags:             tags,
		Favorite:         fromBool(d.Favorite),
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
	if s == "" {
		return nil, nil
	}

	r := csv.NewReader(strings.NewReader(s))
	return r.Read()
}

func fromBool(v bool) int {
	if v {
		return 1
	} else {
		return 0
	}
}

func toBool(i int) bool {
	return i == 1
}

func setDrinkDataAttributes(obj DrinkDataSetter, data DrinkDataGetter) {
	obj.SetName(data.GetName())
	obj.SetPrimaryAlcohol(data.GetPrimaryAlcohol())
	obj.SetPreferredGlass(data.GetPreferredGlass())
	obj.SetIngredients(data.GetIngredients())
	obj.SetInstructions(data.GetInstructions())
	obj.SetNotes(data.GetNotes())
	obj.SetPublicity(data.GetPublicity())
	obj.SetUnderDevelopment(data.GetUnderDevelopment())
	obj.SetTags(data.GetTags())
	obj.SetFavorite(data.GetFavorite())
}
