package drink

import (
	"database/sql"

	"github.com/huandu/go-sqlbuilder"
	"github.com/nicjohnson145/mixer-service/pkg/common"
)

var ModelStruct = sqlbuilder.NewStruct(new(Model))

const (
	TableName = "drink"
)

type Model struct {
	ID               int64  `db:"id"`
	Name             string `db:"name" fieldtag:"required_insert"`
	Username         string `db:"username" fieldtag:"required_insert"`
	PrimaryAlcohol   string `db:"primary_alcohol" fieldtag:"required_insert"`
	PreferredGlass   string `db:"preferred_glass" fieldtag:"required_insert"`
	Ingredients      string `db:"ingredients" fieldtag:"required_insert"`
	Instructions     string `db:"instructions" fieldtag:"required_insert"`
	Notes            string `db:"notes" fieldtag:"required_insert"`
	Publicity        string `db:"publicity" fieldtag:"required_insert"`
	UnderDevelopment int    `db:"under_development" fieldtag:"required_insert"`
	Tags             string `db:"tags" fieldtag:"required_insert"`
	Favorite         int    `db:"favorite" fieldtag:"required_insert"`
}

func getByID(id int64, db *sql.DB) (*Model, error) {
	sb := ModelStruct.SelectFrom(TableName)
	sb.Where(sb.Equal("id", id))

	sql, args := sb.Build()
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hasRow := rows.Next()
	if !hasRow {
		return nil, common.ErrNotFound
	}

	var drink Model
	err = rows.Scan(ModelStruct.Addr(&drink)...)
	if err != nil {
		return nil, err
	}

	return &drink, nil
}

func getByNameAndUsername(name string, username string, db *sql.DB) (*Model, error) {
	sb := ModelStruct.SelectFrom(TableName)
	sb.Where(
		sb.Equal("name", name),
		sb.Equal("username", username),
	)

	sql, args := sb.Build()
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hasRow := rows.Next()
	if !hasRow {
		return nil, common.ErrNotFound
	}

	var drink Model
	err = rows.Scan(ModelStruct.Addr(&drink)...)
	if err != nil {
		return nil, err
	}

	return &drink, nil
}

func getAllPublicDrinksByUser(username string, db *sql.DB) ([]Model, error) {
	drinks, err := getAllDrinksByUser(username, db)
	if err != nil {
		return []Model{}, err
	}

	publicDrinks := make([]Model, 0, len(drinks))
	for _, d := range drinks {
		if d.Publicity == DrinkPublicityPublic {
			publicDrinks = append(publicDrinks, d)
		}
	}
	return publicDrinks, nil
}

func getAllDrinksByUser(username string, db *sql.DB) ([]Model, error) {
	sb := ModelStruct.SelectFrom(TableName)
	sb.Where(
		sb.Equal("username", username),
	)
	sb.OrderBy("id")
	sb.Asc()

	sql, args := sb.Build()
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	drinks := []Model{}
	for rows.Next() {
		var d Model
		err := rows.Scan(ModelStruct.Addr(&d)...)
		if err != nil {
			return []Model{}, err
		}
		drinks = append(drinks, d)
	}

	return drinks, rows.Err()
}

func updateModel(model Model, db *sql.DB) error {
	sb := ModelStruct.Update(TableName, model)
	sb.Where(sb.Equal("id", model.ID))
	sql, args := sb.Build()
	_, err := db.Exec(sql, args...)
	return err
}

func deleteModel(id int64, db *sql.DB) error {
	sb := ModelStruct.DeleteFrom(TableName)
	sb.Where(sb.Equal("id", id))
	sql, args := sb.Build()
	_, err := db.Exec(sql, args...)
	return err
}

func create(d Model, db *sql.DB) (int64, error) {
	sql, args := ModelStruct.InsertIntoForTag(TableName, "required_insert", d).Build()
	rows, err := db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}
