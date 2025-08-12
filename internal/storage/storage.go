package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type ShortURL struct {
	ID          string `json:"id"`
	Original    string `json:"original"`
	AccessCount int    `json:"access_count"`
}

type Storage interface {
	Init() error
	CreateShortURL(id, originalURL string) (int, error)
	GetOriginalURL(id string) (string, error)
	IncrementAccessCount(id string) error
	Ping() error
	GetShortURL(id string) (*ShortURL, error)
}

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage() *SQLiteStorage {
	return &SQLiteStorage{}
}

func (s *SQLiteStorage) Init() error {
	log.Default().Println("Initializing database...")

	var err error
	s.db, err = sql.Open("sqlite3", "./links.db")
	if err != nil {
		return err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS urls (
		id TEXT PRIMARY KEY,
		original_url TEXT NOT NULL,
		access_count INTEGER DEFAULT 0
	);`

	_, err = s.db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Default().Println("Database initialized successfully.")
	return nil
}

func (s *SQLiteStorage) CreateShortURL(id, originalURL string) (int, error) {
	stmt, err := s.db.Prepare("INSERT INTO urls(id, original_url) VALUES(?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, originalURL)
	if err != nil {
		return 0, err
	}
	return 0, nil // Initial access count is 0
}

func (s *SQLiteStorage) GetOriginalURL(id string) (string, error) {
	var originalURL string
	err := s.db.QueryRow("SELECT original_url FROM urls WHERE id = ?", id).Scan(&originalURL)
	return originalURL, err
}

func (s *SQLiteStorage) IncrementAccessCount(id string) error {
	_, err := s.db.Exec("UPDATE urls SET access_count = access_count + 1 WHERE id = ?", id)
	return err
}

func (s *SQLiteStorage) GetShortURL(id string) (*ShortURL, error) {
	var originalURL string
	var accessCount int
	err := s.db.QueryRow("SELECT original_url, access_count FROM urls WHERE id = ?", id).Scan(&originalURL, &accessCount)
	if err != nil {
		return nil, err
	}
	return &ShortURL{ID: id, Original: originalURL, AccessCount: accessCount}, nil
}

func (s *SQLiteStorage) Ping() error {
	return s.db.Ping()
}
