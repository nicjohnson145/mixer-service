package drink

type DrinkData struct {
	Name             string   `json:"name" validate:"required"`
	PrimaryAlcohol   string   `json:"primary_alcohol" validate:"required"`
	PreferredGlass   string   `json:"preferred_glass,omitempty"`
	Ingredients      []string `json:"ingredients" validate:"required"`
	Instructions     string   `json:"instructions,omitempty"`
	Notes            string   `json:"notes,omitempty"`
	Publicity        string   `json:"publicity" validate:"required"`
	UnderDevelopment bool     `json:"under_development"`
	Tags             []string `json:"tags"`
}

func (d *DrinkData) SetName(v string) {
	d.Name = v
}

func (d DrinkData) GetName() string {
	return d.Name
}

func (d *DrinkData) SetPrimaryAlcohol(v string) {
	d.PrimaryAlcohol = v
}

func (d DrinkData) GetPrimaryAlcohol() string {
	return d.PrimaryAlcohol
}

func (d *DrinkData) SetPreferredGlass(v string) {
	d.PreferredGlass = v
}

func (d DrinkData) GetPreferredGlass() string {
	return d.PreferredGlass
}

func (d *DrinkData) SetIngredients(v []string) {
	d.Ingredients = v
}

func (d DrinkData) GetIngredients() []string {
	return d.Ingredients
}

func (d *DrinkData) SetInstructions(v string) {
	d.Instructions = v
}

func (d DrinkData) GetInstructions() string {
	return d.Instructions
}

func (d *DrinkData) SetNotes(v string) {
	d.Notes = v
}

func (d DrinkData) GetNotes() string {
	return d.Notes
}

func (d *DrinkData) SetPublicity(v string) {
	d.Publicity = v
}

func (d DrinkData) GetPublicity() string {
	return d.Publicity
}

func (d *DrinkData) SetUnderDevelopment(v bool) {
	d.UnderDevelopment = v
}

func (d DrinkData) GetUnderDevelopment() bool {
	return d.UnderDevelopment
}

func (d *DrinkData) SetTags(v []string) {
	d.Tags = v
}

func (d DrinkData) GetTags() []string {
	return d.Tags
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
	SetUnderDevelopment(bool)
	SetTags([]string)
}

type DrinkDataGetter interface {
	GetName() string
	GetPrimaryAlcohol() string
	GetPreferredGlass() string
	GetIngredients() []string
	GetInstructions() string
	GetNotes() string
	GetPublicity() string
	GetUnderDevelopment() bool
	GetTags() []string
}
