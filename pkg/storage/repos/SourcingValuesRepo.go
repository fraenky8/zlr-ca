package repos

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/domain"
	"github.com/fraenky8/zlr-ca/pkg/storage"
	"github.com/fraenky8/zlr-ca/pkg/storage/dtos"
)

type SourcingValuesRepo struct {
	db *storage.Database
}

func NewSourcingValuesRepo(db *storage.Database) *SourcingValuesRepo {
	return &SourcingValuesRepo{
		db: db,
	}
}

func (r *SourcingValuesRepo) Creates(sourcingValues domain.SourcingValues) ([]int, error) {

	stmt, err := r.db.Preparex(fmt.Sprintf(`
		INSERT INTO %s.sourcing_values (description) VALUES (TRIM($1)) 
		ON CONFLICT (description) DO UPDATE SET description = TRIM($1) RETURNING id
	`, r.db.Schema))

	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}

	var ids []int
	for _, sourcingValue := range sourcingValues {
		var id int
		err := stmt.Get(&id, sourcingValue)
		if err != nil {
			return nil, fmt.Errorf("could not create sourcing value: %v", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *SourcingValuesRepo) Read(icecreamProductId int) (domain.SourcingValues, error) {

	var sourcingValues []*dtos.SourcingValues
	err := r.db.Select(&sourcingValues, fmt.Sprintf(`
		SELECT
  			id, description
		FROM
  			%s.sourcing_values AS sv,
  			%s.icecream_has_sourcing_values AS ihsv
		WHERE ihsv.sourcing_values_id = sv.id
		AND ihsv.icecream_product_id = $1
	`, r.db.Schema, r.db.Schema), icecreamProductId)

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
	err := r.db.Select(&sourcingValues, fmt.Sprintf(`
		SELECT id, description
		FROM %s.sourcing_values
	`, r.db.Schema))

	if err != nil {
		return nil, err
	}

	return r.convert(sourcingValues)
}

func (r *SourcingValuesRepo) Deletes(icecreamProductIds []int) (err error) {

	tx := r.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	stmt, err := r.db.Preparex(fmt.Sprintf(`
		DELETE FROM %s.icecream_has_sourcing_values
		WHERE icecream_product_id = $1
	`, r.db.Schema))

	if err != nil {
		return fmt.Errorf("could not prepare statement: %v", err)
	}

	for _, id := range icecreamProductIds {
		if _, err := stmt.Exec(id); err != nil {
			return fmt.Errorf("could not delete sourcing values of icecream with productID = %d: %v", id, err)
		}
	}

	return nil
}

func (r *SourcingValuesRepo) convert(sourcingValues []*dtos.SourcingValues) (domain.SourcingValues, error) {
	sv := domain.SourcingValues{}
	for _, i := range sourcingValues {
		sv = append(sv, domain.SourcingValue(i.Description))
	}
	return sv, nil
}
