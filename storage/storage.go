package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type Config struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	MCU       string `json:"mcu"`
	Settings  string `json:"settings"`
}

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) (*Store, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}

	schema := `
	CREATE TABLE IF NOT EXISTS configurations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		mcu TEXT NOT NULL,
		settings TEXT NOT NULL
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Save(cfg Config) (int64, error) {
	if cfg.Name == "" {
		return 0, errors.New("configuration name is required")
	}
	if cfg.MCU == "" {
		return 0, errors.New("MCU model is required")
	}

	if !json.Valid([]byte(cfg.Settings)) {
		return 0, errors.New("settings must contain valid JSON")
	}

	result, err := s.db.Exec(
		`INSERT INTO configurations (name, mcu, settings)
		 VALUES (?, ?, ?)`,
		cfg.Name, cfg.MCU, cfg.Settings,
	)
	if err != nil {
		return 0, fmt.Errorf("save configuration: %w", err)
	}

	return result.LastInsertId()
}

func (s *Store) Get(id int64) (Config, error) {
	var cfg Config

	err := s.db.QueryRow(
		`SELECT id, name, mcu, settings
		 FROM configurations WHERE id = ?`,
		id,
	).Scan(&cfg.ID, &cfg.Name, &cfg.MCU, &cfg.Settings)

	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
