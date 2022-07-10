package settings

import (
	"database/sql"
	"strconv"

	"github.com/huandu/go-sqlbuilder"
)

var ModelStruct = sqlbuilder.NewStruct(new(Model))

const (
	TableName = "user_setting"
)

type Model struct {
	Username string `db:"username" fieldtag:"required_insert"`
	Key      string `db:"key" fieldtag:"required_insert"`
	Value    string `db:"value" fieldtag:"required_insert"`
}

func getByUsername(username string, db *sql.DB) (map[string]string, error) {
	sb := ModelStruct.SelectFrom(TableName)
	sb.Where(sb.Equal("username", username))

	sql, args := sb.Build()
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allSettings := map[string]string{}

	var m Model
	for rows.Next() {
		err := rows.Scan(ModelStruct.Addr(&m)...)
		if err != nil {
			return nil, err
		}
		allSettings[m.Key] = m.Value
	}
	return allSettings, nil
}


func writeSettingsForUser(db *sql.DB, username string, settings UserSettings) error {
	ib := sqlbuilder.NewInsertBuilder()
	ib.ReplaceInto(TableName)
	ib.Cols("username", "key", "value")

	// One rew per setting to persist
	ib.Values(username, PublicProfile, strconv.FormatBool(settings.PublicProfile))

	sql, args := ib.Build()
	_, err := db.Exec(sql, args...)
	return err
}
