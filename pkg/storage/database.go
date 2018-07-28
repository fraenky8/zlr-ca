package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Schema   string
}

type Database interface {
	Connect() error
	Close() error
	DB() *sqlx.DB
	Config() *Config
}

type Postgres struct {
	db  *sqlx.DB
	cfg *Config
}

func NewPostgres(config *Config) (*Postgres, error) {
	if config.Port == "" {
		config.Port = "5432"
	}

	pg := &Postgres{cfg: config}

	if err := pg.Connect(); err != nil {
		return nil, err
	}

	return pg, nil
}

func (pg *Postgres) Connect() (err error) {
	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		pg.cfg.Host, pg.cfg.Port, pg.cfg.Username, pg.cfg.Database, pg.cfg.Password)

	pg.db, err = sqlx.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}

	if err = pg.db.Ping(); err != nil {
		return fmt.Errorf("database not reachable: %v", err)
	}

	return nil
}

func (pg *Postgres) Close() error {
	return pg.db.Close()
}

func (pg *Postgres) DB() *sqlx.DB {
	return pg.db
}

func (pg *Postgres) Config() *Config {
	return pg.cfg
}
