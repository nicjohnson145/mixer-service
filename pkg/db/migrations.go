package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

func createMigrations() *migrate.MemoryMigrationSource {
	idGen := func() func() string {
		n := 0
		return func() string {
			n += 1
			return fmt.Sprint(n)
		}
	}

	nextId := idGen()

	return &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: nextId(),
				Up: []string{
					`
					CREATE TABLE user (
						username TEXT PRIMARY KEY,
						password TEXT NOT NULL
					);
					`,
				},
				Down: []string{"DROP TABLE user;"},
			},
		},
	}
}

func NewDBOrDie(path string) *sql.DB {
	db, err := NewDB(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("opening & migrating db: %v", err))
	}

	return db
}

func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = migrate.Exec(db, "sqlite3", createMigrations(), migrate.Up)
	if err != nil {
		return nil, err
	}

	return db, nil
}