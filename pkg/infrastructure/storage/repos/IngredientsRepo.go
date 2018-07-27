package repos

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/core/domain"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage/dtos"
	"github.com/jmoiron/sqlx"
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

	stmt, err := r.prepareCreateStmt()
	if err != nil {
		return 0, err
	}

	return r.create(stmt, ingredient)
}

func (r *IngredientsRepo) Creates(ingredients domain.Ingredients) ([]int64, error) {

	stmt, err := r.prepareCreateStmt()
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, ingredient := range ingredients {

		id, err := r.create(stmt, ingredient)
		if err != nil {
			return nil, fmt.Errorf("could not create ingredients: %v", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *IngredientsRepo) prepareCreateStmt() (*sqlx.Stmt, error) {
	stmt, err := r.db.Preparex(`
		INSERT INTO ingredients (name) VALUES (TRIM($1)) 
		ON CONFLICT (name) DO UPDATE SET name = TRIM($1) RETURNING id
	`)
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}
	return stmt, nil
}

func (r *IngredientsRepo) create(stmt *sqlx.Stmt, ingredient domain.Ingredient) (int64, error) {
	var id int64
	err := stmt.Get(&id, ingredient)
	if err != nil {
		return 0, fmt.Errorf("could not create ingredient: %v", err)
	}
	return id, nil
}

func (r *IngredientsRepo) Read(icecreamProductId int64) (domain.Ingredients, error) {

	var ingredients []*dtos.Ingredients
	err := r.db.Select(&ingredients, `
		SELECT
  			id, name
		FROM
  			ingredients AS i,
  			icecream_has_ingredients AS ihi
		WHERE ihi.ingredients_id = i.id
		AND ihi.icecream_product_id = $1
	`, icecreamProductId)

	if err != nil {
		return nil, err
	}

	return r.Convert(ingredients)
}

func (r *IngredientsRepo) Convert(ingredients []*dtos.Ingredients) (domain.Ingredients, error) {
	di := domain.Ingredients{}
	for _, i := range ingredients {
		di = append(di, domain.Ingredient(i.Name))
	}
	return di, nil
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
