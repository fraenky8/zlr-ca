package repos

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/core/domain"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage/dtos"
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

func (r *SourcingValuesRepo) Create(sourcingValue domain.SourcingValue) (int, error) {

	stmt, err := r.prepareCreateStmt()
	if err != nil {
		return 0, err
	}

	return r.create(stmt, sourcingValue)
}

func (r *SourcingValuesRepo) Creates(sourcingValues domain.SourcingValues) ([]int, error) {

	stmt, err := r.prepareCreateStmt()
	if err != nil {
		return nil, err
	}

	var ids []int
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

func (r *SourcingValuesRepo) create(stmt *sqlx.Stmt, sourcingValue domain.SourcingValue) (int, error) {
	var id int
	err := stmt.Get(&id, sourcingValue)
	if err != nil {
		return 0, fmt.Errorf("could not create sourcing value: %v", err)
	}
	return id, nil
}

func (r *SourcingValuesRepo) Read(icecreamProductId int) (domain.SourcingValues, error) {

	var sourcingValues []*dtos.SourcingValues
	err := r.db.Select(&sourcingValues, `
		SELECT
  			id, description
		FROM
  			sourcing_values AS sv,
  			icecream_has_sourcing_values AS ihsv
		WHERE ihsv.sourcing_values_id = sv.id
		AND ihsv.icecream_product_id = $1
	`, icecreamProductId)

	if err != nil {
		return nil, err
	}

	return r.convert(sourcingValues)
}

func (r *SourcingValuesRepo) Reads(icecreamProductIds []int) (sourcingValues []domain.SourcingValues, err error) {
	for _, id := range icecreamProductIds {
		sourcingValue, err := r.Read(id)
		if err != nil {
			return nil, err
		}
		sourcingValues = append(sourcingValues, sourcingValue)
	}
	return sourcingValues, nil
}

func (r *SourcingValuesRepo) ReadAll() (domain.SourcingValues, error) {

	var sourcingValues []*dtos.SourcingValues
	err := r.db.Select(&sourcingValues, `
		SELECT id, description
		FROM sourcing_values
	`)

	if err != nil {
		return nil, err
	}

	return r.convert(sourcingValues)
}

func (r *SourcingValuesRepo) convert(sourcingValues []*dtos.SourcingValues) (domain.SourcingValues, error) {
	sv := domain.SourcingValues{}
	for _, i := range sourcingValues {
		sv = append(sv, domain.SourcingValue(i.Description))
	}
	return sv, nil
}
