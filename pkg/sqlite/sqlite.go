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
	sqlite           *sql.DB
}

// Логику sqlite мы вымещаем вне internal т.к это не критичная и переиспользуемая логика.

// Connecttion
func (d *Database) ConnectionToDB(cfg *config.Config) (db *Database, err error) {
	if cfg == nil {
		return nil, errors.New("error, config is nil")
	}

	db.sqlite, err = sql.Open("go-sqlite3", cfg.DBConfig.DatabaseFilePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Disconnect
func (d *Database) Disconnect() error {
	if err := d.sqlite.Close(); err != nil {
		return err
	}

	return nil
}
