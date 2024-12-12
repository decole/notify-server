package postgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq" // init sql driver
	"notify-server/internal/config"
)

type Storage struct {
	db *sql.DB
}

func New(storage config.Storage) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		storage.Host, storage.Port, storage.User, storage.Password, storage.DatabaseName)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return &Storage{db: db}, nil
}

func (s Storage) SaveNotify(client string, message string) error {
	const op = "storage.postgres.SaveNotify"

	stmt, err := s.db.Prepare("INSERT INTO notify(client, message) VALUES($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	_, err = stmt.Exec(client, message)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			return fmt.Errorf("%s: %s", op, pqErr)
		}

		return fmt.Errorf("%s: %s", op, err)
	}

	return err
}

func (s Storage) GetNotify(client string) (string, error) {
	const op = "storage.postgres.GetNotify"

	stmt, err := s.db.Prepare("SELECT message FROM notify WHERE client = $1")
	if err != nil {
		return "", fmt.Errorf("%s: %s", op, err)
	}

	var message string
	err = stmt.QueryRow(client).Scan(&message)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			return "", fmt.Errorf("%s: %s", op, pqErr)
		}

		return "", fmt.Errorf("%s: execcute statement: %s", op, err)
	}

	return message, err
}
