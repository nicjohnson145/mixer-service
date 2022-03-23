package drink

type drinkData struct {
	Name           string   `json:"name" validate:"required"`
	PrimaryAlcohol string   `json:"primary_alcohol" validate:"required"`
	PreferredGlass string   `json:"preferred_glass"`
	Ingredients    []string `json:"ingredients" validate:"required"`
	Instructions   string   `json:"instructions"`
	Notes          string   `json:"notes"`
	Publicity      string   `json:"publicity" validate:"required"`
}

func (d *drinkData) SetName(v string) {
	d.Name = v
}

func (d drinkData) GetName() string {
	return d.Name
}

func (d *drinkData) SetPrimaryAlcohol(v string) {
	d.PrimaryAlcohol = v
}

func (d drinkData) GetPrimaryAlcohol() string {
	return d.PrimaryAlcohol
}

func (d *drinkData) SetPreferredGlass(v string) {
	d.PreferredGlass = v
}

func (d drinkData) GetPreferredGlass() string {
	return d.PreferredGlass
}

func (d *drinkData) SetIngredients(v []string) {
	d.Ingredients = v
}

func (d drinkData) GetIngredients() []string {
	return d.Ingredients
}

func (d *drinkData) SetInstructions(v string) {
	d.Instructions = v
}

func (d drinkData) GetInstructions() string {
	return d.Instructions
}

func (d *drinkData) SetNotes(v string) {
	d.Notes = v
}

func (d drinkData) GetNotes() string {
	return d.Notes
}

func (d *drinkData) SetPublicity(v string) {
	d.Publicity = v
}

func (d drinkData) GetPublicity() string {
	return d.Publicity
}

type DrinkDataOperator interface {
	DrinkDataSetter
	DrinkDataGetter
}

type DrinkDataSetter interface {
	SetName(string)
	SetPrimaryAlcohol(string)
	SetPreferredGlass(string)
	SetIngredients([]string)
	SetInstructions(string)
	SetNotes(string)
	SetPublicity(string)
}

type DrinkDataGetter interface {
	GetName() string
	GetPrimaryAlcohol() string
	GetPreferredGlass() string
	GetIngredients() []string
	GetInstructions() string
	GetNotes() string
	GetPublicity() string
}
