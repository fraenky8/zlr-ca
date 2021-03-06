package repos

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/domain"
	"github.com/fraenky8/zlr-ca/pkg/storage"
	"github.com/fraenky8/zlr-ca/pkg/storage/dtos"
)

type IngredientsRepo struct {
	db storage.Database
}

func NewIngredientsRepo(db storage.Database) *IngredientsRepo {
	return &IngredientsRepo{
		db: db,
	}
}

func (r *IngredientsRepo) Creates(ingredients domain.Ingredients) ([]int64, error) {

	stmt, err := r.db.DB().Preparex(fmt.Sprintf(`
		INSERT INTO %s.ingredients (name) VALUES (TRIM($1)) 
		ON CONFLICT (name) DO UPDATE SET name = TRIM($1) RETURNING id
	`, r.db.Config().Schema))

	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}

	var ids []int64
	for _, ingredient := range ingredients {
		var id int64
		err := stmt.Get(&id, ingredient)
		if err != nil {
			return nil, fmt.Errorf("could not create ingredient: %v", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *IngredientsRepo) Read(icecreamProductId int64) (domain.Ingredients, error) {

	var ingredients []*dtos.Ingredients
	err := r.db.DB().Select(&ingredients, fmt.Sprintf(`
		SELECT
  			id, name
		FROM
  			%s.ingredients AS i,
  			%s.icecream_has_ingredients AS ihi
		WHERE ihi.ingredients_id = i.id
		AND ihi.icecream_product_id = $1
	`, r.db.Config().Schema, r.db.Config().Schema), icecreamProductId)

	if err != nil {
		return nil, err
	}

	return r.convert(ingredients)
}

func (r *IngredientsRepo) Reads(icecreamProductIds []int64) (ingredients []domain.Ingredients, err error) {
	for _, id := range icecreamProductIds {
		ingredient, err := r.Read(id)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

func (r *IngredientsRepo) ReadAll() (domain.Ingredients, error) {

	var ingredients []*dtos.Ingredients
	err := r.db.DB().Select(&ingredients, fmt.Sprintf(`
		SELECT id, name
		FROM %s.ingredients
	`, r.db.Config().Schema))

	if err != nil {
		return nil, err
	}

	return r.convert(ingredients)
}

func (r *IngredientsRepo) convert(ingredients []*dtos.Ingredients) (domain.Ingredients, error) {
	di := domain.Ingredients{}
	for _, i := range ingredients {
		di = append(di, domain.Ingredient(i.Name))
	}
	return di, nil
}
