package repos

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/storage"
)

type IcecreamHasSourcingValuesRepo struct {
	db storage.Database
}

func NewIcecreamHasSourcingValuesRepo(db storage.Database) *IcecreamHasSourcingValuesRepo {
	return &IcecreamHasSourcingValuesRepo{
		db: db,
	}
}

func (r *IcecreamHasSourcingValuesRepo) Create(productId int64, sourcingValueIds []int64) error {
	stmt, err := r.db.DB().Preparex(fmt.Sprintf(`
		INSERT INTO %s.icecream_has_sourcing_values 
			(icecream_product_id, sourcing_values_id) 
		VALUES ($1, $2) 
		ON CONFLICT (icecream_product_id, sourcing_values_id) DO NOTHING
	`, r.db.Config().Schema))

	if err != nil {
		return fmt.Errorf("could not prepare statement: %v", err)
	}

	for _, id := range sourcingValueIds {
		if _, err = stmt.Exec(productId, id); err != nil {
			return fmt.Errorf("could not create sourcing value relationship: %v", err)
		}
	}

	return nil
}
