package repos

import (
	"fmt"
	"github.com/fraenky8/zlr-ca/pkg/core/domain"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
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

	stmt, err := r.db.Preparex(`
		INSERT INTO sourcing_values (description) VALUES (TRIM($1)) 
		ON CONFLICT (description) DO UPDATE SET description = TRIM($1) RETURNING id
	`)

	if err != nil {
		return 0, fmt.Errorf("could not prepare statement: %v", err)
	}

	var id int64
	err = stmt.Get(&id, sourcingValue)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve last inserted id: %v", err)
	}

	return id, nil
}

func (r *SourcingValuesRepo) Creates(sourcingValues domain.SourcingValues) ([]int64, error) {

	stmt, err := r.db.Preparex(`
		INSERT INTO sourcing_values (description) VALUES (TRIM($1)) 
		ON CONFLICT (description) DO UPDATE SET description = TRIM($1) RETURNING id
	`)

	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}

	var ids []int64
	for _, sourcingValue := range sourcingValues {

		var id int64
		err = stmt.Get(&id, sourcingValue)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve last inserted id: %v", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
