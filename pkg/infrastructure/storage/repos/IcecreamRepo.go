package repos

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/core/domain"
	"github.com/fraenky8/zlr-ca/pkg/infrastructure/storage"
	"github.com/jmoiron/sqlx"
)

type IcecreamRepo struct {
	db *storage.Database
}

func NewIcecreamRepo(db *storage.Database) *IcecreamRepo {
	return &IcecreamRepo{
		db: db,
	}
}

func (r *IcecreamRepo) Create(ic domain.Icecream) (int64, error) {

	stmt, err := r.prepareCreateStmt()
	if err != nil {
		return 0, err
	}

	return r.create(stmt, ic)
}

func (r *IcecreamRepo) Creates(icecreams []domain.Icecream) ([]int64, error) {

	stmt, err := r.prepareCreateStmt()
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, icecream := range icecreams {

		id, err := r.create(stmt, icecream)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *IcecreamRepo) prepareCreateStmt() (*sqlx.Stmt, error) {
	stmt, err := r.db.Preparex(`
		INSERT INTO icecream
  			(name, description, story, image_open, image_closed, allergy_info, dietary_certifications)
		VALUES
  			($1, $2, $3, $4, $5, $6, $7)
		RETURNING product_id
	`)
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}
	return stmt, nil
}

func (r *IcecreamRepo) create(stmt *sqlx.Stmt, ic domain.Icecream) (int64, error) {

	var productId int64
	err := stmt.Get(&productId, ic.Name, ic.Description, ic.Story, ic.ImageOpen, ic.ImageClosed, ic.AllergyInfo, ic.DietaryCertifications)
	if err != nil {
		return 0, fmt.Errorf("could not create icecream: %v", err)
	}

	ids, err := NewIngredientsRepo(r.db).Creates(ic.Ingredients)
	if err != nil {
		return 0, fmt.Errorf("could not create ingredients: %v", err)
	}

	err = NewIcecreamHasIngredientsRepo(r.db).Create(productId, ids)
	if err != nil {
		return 0, fmt.Errorf("could not create ingredients relationships: %v", err)
	}

	ids, err = NewSourcingValuesRepo(r.db).Creates(ic.SourcingValues)
	if err != nil {
		return 0, fmt.Errorf("could not create sourcing values: %v", err)
	}

	err = NewIcecreamHasSourcingValuesRepo(r.db).Create(productId, ids)
	if err != nil {
		return 0, fmt.Errorf("could not create sourcing values relationships: %v", err)
	}

	return productId, nil
}
