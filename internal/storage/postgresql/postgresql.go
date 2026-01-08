package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(connStr string) (*Storage, error) {
	const op = "storage.postgresql.new"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS url(
        id SERIAL PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    ) 
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
func (s Storage) Save(urlToSave, alias string) error {
	const op = "storage.postgresql.New"

	_, err := s.db.Exec(
		`INSERT INTO url(url, alias) VALUES ($1, $2)
     ON CONFLICT (alias) DO UPDATE SET url = excluded.url`,
		urlToSave, alias,
	)

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}
func (s Storage) Get(alias string) (string, error) {
	const op = "storage.postgresql.Get"
	urlToFind := ""
	err := s.db.QueryRow(`select url from url where alias = $1`, alias).Scan(&urlToFind)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("no rows in result by alias %s", alias)
	}
	if err != nil {
		return "", fmt.Errorf("%s,%w", op, err)
	}

	return urlToFind, nil
}
func (s Storage) Delete(alias string) error {
	const op = "storage.postgresql.Delete"
	_, err := s.db.Exec(`delete from url where alias = $1`, alias)
	if err != nil {
		return fmt.Errorf("%s,%w", op, err)
	}
	return nil
}
func (s *Storage) Update(alias string, newURL string) error {
	const op = "storage.postgresql.Update"

	_, err := s.db.Exec(
		`UPDATE url SET url = $1 WHERE alias = $2`, newURL, alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (s *Storage) Close() error {
	const op = "storage.postgresql.Close"

	if s.db == nil {
		return nil
	}

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.db = nil
	return nil
}
