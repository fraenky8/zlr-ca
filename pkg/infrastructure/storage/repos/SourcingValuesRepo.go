package repos

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/core/domain"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
	"github.com/jmoiron/sqlx"
)

type SourcingValuesRepo struct {
	db *storage.Database
}

func NewSourcingValuesRepo(db *storage.Database) *SourcingValuesRepo {
	return &SourcingValuesRepo{
		db: db,
	}
}

func (r *SourcingValuesRepo) Create(sourcingValue domain.SourcingValue) (int64, error) {

	stmt, err := r.prepareCreateStmt()
	if err != nil {
		return 0, err
	}

	return r.create(stmt, sourcingValue)
}

func (r *SourcingValuesRepo) Creates(sourcingValues domain.SourcingValues) ([]int64, error) {

	stmt, err := r.prepareCreateStmt()
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, sourcingValue := range sourcingValues {

		id, err := r.create(stmt, sourcingValue)
		if err != nil {
			return nil, fmt.Errorf("could not create sourcing values: %v", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *SourcingValuesRepo) prepareCreateStmt() (*sqlx.Stmt, error) {
	stmt, err := r.db.Preparex(`
		INSERT INTO sourcing_values (description) VALUES (TRIM($1)) 
		ON CONFLICT (description) DO UPDATE SET description = TRIM($1) RETURNING id
	`)
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}
	return stmt, nil
}

func (r *SourcingValuesRepo) create(stmt *sqlx.Stmt, sourcingValue domain.SourcingValue) (int64, error) {
	var id int64
	err := stmt.Get(&id, sourcingValue)
	if err != nil {
		return 0, fmt.Errorf("could not create sourcing value: %v", err)
	}
	return id, nil
}
