package drink

//go:generate strugen -t DrinkData gen-template.txt
type DrinkData struct {
	Name             string   `json:"name" validate:"required" strugen:"read,write"`
	PrimaryAlcohol   string   `json:"primary_alcohol" validate:"required" strugen:"read,write"`
	PreferredGlass   string   `json:"preferred_glass,omitempty" strugen:"read,write"`
	Ingredients      []string `json:"ingredients" validate:"required" strugen:"read,write"`
	Instructions     string   `json:"instructions,omitempty" strugen:"read,write"`
	Notes            string   `json:"notes,omitempty" strugen:"read,write"`
	Publicity        string   `json:"publicity" validate:"required" strugen:"read,write"`
	UnderDevelopment bool     `json:"under_development" strugen:"read,write"`
	Tags             []string `json:"tags" strugen:"read,write"`
	Favorite         bool     `json:"favorite" strugen:"read,write"`
}

type DrinkDataOperator interface {
	DrinkDataSetter
	DrinkDataGetter
}
