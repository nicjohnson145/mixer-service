package db

import (
	"database/sql"
	"fmt"
	"github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
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
			{
				Id: nextId(),
				Up: []string{
					`
					CREATE TABLE drink (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT NOT NULL,
						username TEXT NOT NULL,
						primary_alcohol TEXT NOT NULL,
						preferred_glass TEXT,
						ingredients TEXT NOT NULL,
						instructions TEXT,
						notes TEXT,
						UNIQUE(name, username),
						FOREIGN KEY (username) REFERENCES user(username) ON DELETE CASCADE ON UPDATE CASCADE
					);
					`,
				},
				Down: []string{
					`
					DROP TABLE drink;
					`,
				},
			},
			{
				Id: nextId(),
				Up: []string{
					`
					ALTER TABLE drink ADD COLUMN publicity TEXT NOT NULL;
					`,
				},
				Down: []string{
					`
					ALTER TABLE drink DROP COLUMN publicity;
					`,
				},
			},
			{
				Id: nextId(),
				Up: []string{`
					CREATE TABLE user_setting (
						username TEXT NOT NULL,
						key TEXT NOT NULL,
						value TEXT NOT NULL,
						UNIQUE(username, key)
					);
				`},
				Down: []string{"DROP TABLE user_setting;"},
			},
			{
				Id: nextId(),
				Up: []string{`
					ALTER TABLE
						drink
					ADD COLUMN
						under_development
						INTEGER
					NOT NULL
					DEFAULT
						FALSE
					;
				`},
				Down: []string{`
					ALTER TABLE
						drink
					DROP COLUMN
						under_development
					;
				`},
			},
			{
				Id: nextId(),
				Up: []string{`
					ALTER TABLE
						drink
					ADD COLUMN
						tags TEXT
					;
				`},
				Down: []string{`
					ALTER TABLE
						drink
					DROP COLUMN
						tags
					;
				`},
			},
			//{
			//    Id: nextId(),
			//    Up: []string{},
			//    Down: []string{},
			//},
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
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = migrate.Exec(db, "sqlite3", createMigrations(), migrate.Up)
	if err != nil {
		return nil, err
	}

	return db, nil
}
