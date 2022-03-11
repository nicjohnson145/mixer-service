package auth

import (
	"database/sql"
	"errors"
	"github.com/huandu/go-sqlbuilder"
)

type UserModel struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

const (
	UserTable = "user"
)

var ErrNotFound = errors.New("not found")

var UserModelStruct = sqlbuilder.NewStruct(new(UserModel))

func getUserByName(username string, db *sql.DB) (*UserModel, error) {
	b := UserModelStruct.SelectFrom(UserTable)
	b.Where(b.Equal("username", username))
	sql, args := b.Build()

	rows, err := db.Query(sql, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	hasRow := rows.Next()
	if !hasRow {
		return nil, ErrNotFound
	}

	var user UserModel
	err = rows.Scan(UserModelStruct.Addr(&user)...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func createUser(user UserModel, db *sql.DB) error {
	ib := UserModelStruct.InsertInto(UserTable, user)
	sql, args := ib.Build()
	_, err := db.Exec(sql, args...)
	return err
}