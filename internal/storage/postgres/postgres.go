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

type Client string

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

	stmt, err := s.db.Prepare("SELECT id, message FROM notify WHERE client = $1 AND read_at IS NULL ORDER BY create_at")
	if err != nil {
		return "", fmt.Errorf("%s: %s", op, err)
	}

	var message string
	var id int64
	err = stmt.QueryRow(client).Scan(&id, &message)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			return "", fmt.Errorf("%s: %s", op, pqErr)
		}

		return "", fmt.Errorf("%s: execcute statement: %s", op, err)
	}

	if err == nil {
		stmt, err = s.db.Prepare("UPDATE notify SET read_at = now() WHERE id = $1")
		if err != nil {
			return "", fmt.Errorf("%s: %s", op, err)
		}

		_, err = stmt.Exec(id)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				return "", fmt.Errorf("%s: %s", op, pqErr)
			}

			return "", fmt.Errorf("%s: %s", op, err)
		}
	}

	return message, err
}

func (s Storage) GetActiveUsers() ([]Client, error) {
	rows, err := s.db.Query("SELECT name FROM client WHERE is_active = TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var clients []Client

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var name Client
		if err := rows.Scan(&name); err != nil {
			return clients, err
		}
		clients = append(clients, name)
	}
	if err = rows.Err(); err != nil {
		return clients, err
	}

	return clients, nil
}

func (s Storage) ClientRegistered(client string) (bool, error) {
	const op = "storage.postgres.ClientRegistered"

	stmt, err := s.db.Prepare("SELECT name, is_active FROM client WHERE name = $1")
	if err != nil {
		return false, fmt.Errorf("%s: %s", op, err)
	}

	var name string
	var isActive bool
	err = stmt.QueryRow(client).Scan(&name, &isActive)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			return false, fmt.Errorf("%s: %s", op, pqErr)
		}

		return false, fmt.Errorf("%s: execcute statement: %s", op, err)
	}

	return isActive, err
}

func (s Storage) SaveClient(client string) error {
	const op = "storage.postgres.SaveClient"

	stmt, err := s.db.Prepare("INSERT INTO client(name) VALUES($1)")
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	_, err = stmt.Exec(client)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			return fmt.Errorf("%s: %s", op, pqErr)
		}

		return fmt.Errorf("%s: %s", op, err)
	}

	err = s.SaveNotify(client, "Welcome to notify service")
	if err != nil {
		return fmt.Errorf("%s: %s save welcome message", op, err)
	}

	return err
}
