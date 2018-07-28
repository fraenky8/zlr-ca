package storage

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Service struct {
	Db                               *Database
	IcecreamService                  domain.IcecreamService
	IngredientService                domain.IngredientService
	SourcingValueService             domain.SourcingValueService
	IcecreamHasIngredientsService    domain.IcecreamHasIngredientsService
	IcecreamHasSourcingValuesService domain.IcecreamHasSourcingValuesService
}

func (s *Service) Verify() error {
	if s.Db == nil {
		return fmt.Errorf("no Database given")
	}
	if s.IcecreamService == nil {
		return fmt.Errorf("no IcecreamService given")
	}
	if s.IngredientService == nil {
		return fmt.Errorf("no IngredientService given")
	}
	if s.SourcingValueService == nil {
		return fmt.Errorf("no SourcingValueService given")
	}
	if s.IcecreamHasIngredientsService == nil {
		return fmt.Errorf("no IcecreamHasIngredientsService given")
	}
	if s.IcecreamHasSourcingValuesService == nil {
		return fmt.Errorf("no IcecreamHasSourcingValuesService given")
	}
	return nil
}

type Database struct {
	*sqlx.DB
	*Config
}

type Config struct {
	Host     string
	Username string
	Password string
	Database string
	Schema   string
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
