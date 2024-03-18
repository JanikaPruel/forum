package sqlite

import (
	"database/sql"
	"errors"
	"forum/pkg/config"

	_ "github.com/mattn/go-sqlite3"
)

// sqlite structure
type Database struct {
	// path to database file
	DatabaseFilePath string
	SQLite           *sql.DB
}

// Логику sqlite мы вымещаем вне internal т.к это не критичная и переиспользуемая логика.

// Connecttion
func ConnectionToDB(cfg *config.Config) (*Database, error) {
	if cfg == nil {
		return nil, errors.New("error, config is nil")
	}
	DB := &Database{}
	db, err := sql.Open("sqlite3", cfg.DBConfig.DatabaseFilePath)
	if err != nil {
		return nil, err
	}
	DB.SQLite = db
	return DB, nil
}

// Disconnect
func (d *Database) Disconnect() error {
	if err := d.SQLite.Close(); err != nil {
		return err
	}

	return nil
}
