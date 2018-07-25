package repos

import (
	"fmt"
	"github.com/fraenky8/zlr-ca/pkg/core/domain"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
)

type IngredientsRepo struct {
	db *storage.Database
}

func NewIngredientsRepo(db *storage.Database) *IngredientsRepo {
	return &IngredientsRepo{
		db: db,
	}
}

func (r *IngredientsRepo) Create(ingredient domain.Ingredient) (int64, error) {

	stmt, err := r.db.Preparex(`
		INSERT INTO ingredients (name) VALUES (TRIM($1)) 
		ON CONFLICT (name) DO UPDATE SET name = TRIM($1) RETURNING id
	`)

	if err != nil {
		return 0, fmt.Errorf("could not prepare statement: %v", err)
	}

	var id int64
	err = stmt.Get(&id, ingredient)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve last inserted id: %v", err)
	}

	return id, nil
}

func (r *IngredientsRepo) Creates(ingredients domain.Ingredients) ([]int64, error) {

	stmt, err := r.db.Preparex(`
		INSERT INTO ingredients (name) VALUES (TRIM($1)) 
		ON CONFLICT (name) DO UPDATE SET name = TRIM($1) RETURNING id
	`)

	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}

	var ids []int64
	for _, ingredient := range ingredients {

		var id int64
		err = stmt.Get(&id, ingredient)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve last inserted id: %v", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
