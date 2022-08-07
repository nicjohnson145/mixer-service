package drink

// GENERATED CODE: DO NOT EDIT

// Generated accessor for DrinkData
func (t DrinkData) GetFavorite() bool {
	return t.Favorite
}
func (t *DrinkData) SetFavorite(v bool)  {
	t.Favorite = v
}
func (t DrinkData) GetIngredients() []string {
	return t.Ingredients
}
func (t *DrinkData) SetIngredients(v []string)  {
	t.Ingredients = v
}
func (t DrinkData) GetInstructions() string {
	return t.Instructions
}
func (t *DrinkData) SetInstructions(v string)  {
	t.Instructions = v
}
func (t DrinkData) GetName() string {
	return t.Name
}
func (t *DrinkData) SetName(v string)  {
	t.Name = v
}
func (t DrinkData) GetNotes() string {
	return t.Notes
}
func (t *DrinkData) SetNotes(v string)  {
	t.Notes = v
}
func (t DrinkData) GetPreferredGlass() string {
	return t.PreferredGlass
}
func (t *DrinkData) SetPreferredGlass(v string)  {
	t.PreferredGlass = v
}
func (t DrinkData) GetPrimaryAlcohol() string {
	return t.PrimaryAlcohol
}
func (t *DrinkData) SetPrimaryAlcohol(v string)  {
	t.PrimaryAlcohol = v
}
func (t DrinkData) GetPublicity() string {
	return t.Publicity
}
func (t *DrinkData) SetPublicity(v string)  {
	t.Publicity = v
}
func (t DrinkData) GetTags() []string {
	return t.Tags
}
func (t *DrinkData) SetTags(v []string)  {
	t.Tags = v
}
func (t DrinkData) GetUnderDevelopment() bool {
	return t.UnderDevelopment
}
func (t *DrinkData) SetUnderDevelopment(v bool)  {
	t.UnderDevelopment = v
}

type DrinkDataSetter interface {
	SetFavorite(bool)
	SetIngredients([]string)
	SetInstructions(string)
	SetName(string)
	SetNotes(string)
	SetPreferredGlass(string)
	SetPrimaryAlcohol(string)
	SetPublicity(string)
	SetTags([]string)
	SetUnderDevelopment(bool)
}

type DrinkDataGetter interface {
	GetFavorite() bool
	GetIngredients() []string
	GetInstructions() string
	GetName() string
	GetNotes() string
	GetPreferredGlass() string
	GetPrimaryAlcohol() string
	GetPublicity() string
	GetTags() []string
	GetUnderDevelopment() bool
}

func setDrinkDataAttributes(obj DrinkDataSetter, data DrinkDataGetter) {
	obj.SetFavorite(data.GetFavorite())
	obj.SetIngredients(data.GetIngredients())
	obj.SetInstructions(data.GetInstructions())
	obj.SetName(data.GetName())
	obj.SetNotes(data.GetNotes())
	obj.SetPreferredGlass(data.GetPreferredGlass())
	obj.SetPrimaryAlcohol(data.GetPrimaryAlcohol())
	obj.SetPublicity(data.GetPublicity())
	obj.SetTags(data.GetTags())
	obj.SetUnderDevelopment(data.GetUnderDevelopment())
}


