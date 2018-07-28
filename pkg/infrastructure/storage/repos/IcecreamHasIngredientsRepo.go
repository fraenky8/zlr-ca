package repos

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
)

type IcecreamHasIngredientsRepo struct {
	db *storage.Database
}

func NewIcecreamHasIngredientsRepo(db *storage.Database) *IcecreamHasIngredientsRepo {
	return &IcecreamHasIngredientsRepo{
		db: db,
	}
}

func (r *IcecreamHasIngredientsRepo) Create(productId int, ingredientIds []int) error {
	stmt, err := r.db.Preparex(`
		INSERT INTO icecream_has_ingredients 
			(icecream_product_id, ingredients_id) 
		VALUES ($1, $2)
		ON CONFLICT (icecream_product_id, ingredients_id) DO NOTHING
	`)

	if err != nil {
		return fmt.Errorf("could not prepare statement: %v", err)
	}

	for _, id := range ingredientIds {
		if _, err = stmt.Exec(productId, id); err != nil {
			return fmt.Errorf("could not create ingredient relationship: %v", err)
		}
	}

	return nil
}
