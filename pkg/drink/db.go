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
	ID             int    `db:"id"`
	Name           string `db:"name"`
	Username       string `db:"username"`
	PrimaryAlcohol string `db:"primary_alcohol"`
	PreferredGlass string `db:"preferred_glass"`
	Ingredients    string `db:"ingredients"`
	Instructions   string `db:"instructions"`
	Notes          string `db:"notes"`
}


func getByID(id int, db *sql.DB) (*Model, error) {
	sb := ModelStruct.SelectFrom(TableName)
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

	var drink Model
	err = rows.Scan(ModelStruct.Addr(&drink)...)
	if err != nil {
		return nil, err
	}

	return &drink, nil
}

func create(d Model, db *sql.DB) (int64, error) {
	b := ModelStruct.InsertInto(TableName, d)
	sql, args := b.Build()
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
