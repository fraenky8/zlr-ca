package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Username string
	Password string
	Database string
	Schema   string
}

type Database struct {
	*sqlx.DB
	*Config
}

func Connect(cfg *Config) (*Database, error) {

	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		cfg.Host, "5432", cfg.Username, cfg.Database, cfg.Password)

	sqldb, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	return &Database{sqldb, cfg}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}