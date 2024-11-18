package postgres

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(host string, port int, user string, password string, dbname string) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS notify(
	    id serial CONSTRAINT notify_pk PRIMARY KEY,
	    "user"  VARCHAR(250) NOT NULL,
	    message json NOT NULL
	)
	CREATE INDEX IF NOT EXISTS idx_notify_user ON notify ("user")
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return &Storage{db: db}, nil
}
