package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"community-support-service/internal/config"
)

type DB struct {
	*sqlx.DB
}

func New(cfg *config.Config) (*DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// Transaction helper
func (db *DB) WithTx(fn func(*sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}