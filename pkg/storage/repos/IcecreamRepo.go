package repos

import (
	"fmt"
	"strconv"

	"github.com/fraenky8/zlr-ca/pkg/domain"
	"github.com/fraenky8/zlr-ca/pkg/storage"
	"github.com/fraenky8/zlr-ca/pkg/storage/dtos"
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

func (r *IcecreamRepo) Create(icecream domain.Icecream) (int, error) {

	stmt, err := r.prepareCreateStmt()
	if err != nil {
		return 0, err
	}

	return r.create(stmt, icecream)
}

func (r *IcecreamRepo) Creates(icecreams []domain.Icecream) ([]int, error) {

	stmt, err := r.prepareCreateStmt()
	if err != nil {
		return nil, err
	}

	var ids []int
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
	stmt, err := r.db.Preparex(fmt.Sprintf(`
		INSERT INTO %s.icecream
  			(product_id, name, description, story, image_open, image_closed, allergy_info, dietary_certifications)
		VALUES
  			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING product_id
	`, r.db.Schema))
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}
	return stmt, nil
}

func (r *IcecreamRepo) create(stmt *sqlx.Stmt, icecream domain.Icecream) (int, error) {

	tx := r.db.MustBegin()

	var productId int
	err := stmt.Get(&productId,
		icecream.ProductID, icecream.Name, icecream.Description, icecream.Story,
		icecream.ImageOpen, icecream.ImageClosed, icecream.AllergyInfo, icecream.DietaryCertifications,
	)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("could not create icecream: %v", err)
	}

	ids, err := NewIngredientsRepo(r.db).Creates(icecream.Ingredients)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("could not create ingredients: %v", err)
	}

	err = NewIcecreamHasIngredientsRepo(r.db).Create(productId, ids)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("could not create ingredients relationships: %v", err)
	}

	ids, err = NewSourcingValuesRepo(r.db).Creates(icecream.SourcingValues)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("could not create sourcing values: %v", err)
	}

	err = NewIcecreamHasSourcingValuesRepo(r.db).Create(productId, ids)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("could not create sourcing values relationships: %v", err)
	}

	tx.Commit()

	return productId, nil
}

func (r *IcecreamRepo) Read(ids []int) ([]*domain.Icecream, error) {

	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT 
			product_id, 
			name, 
			description, 
			story, 
			image_open, 
			image_closed, 
			allergy_info, 
			dietary_certifications
		FROM %s.icecream 
		WHERE product_id IN (?)
	`, r.db.Schema), ids)

	if err != nil {
		return nil, err
	}

	// http://jmoiron.github.io/sqlx/#inQueries
	// sqlx.In returns queries with the `?` bindvar, we can rebind it for our backend
	// here: ? to $#
	query = r.db.Rebind(query)

	var icecreamsDtos []dtos.Icecream
	if err = r.db.Select(&icecreamsDtos, query, args...); err != nil {
		return nil, err
	}

	if len(icecreamsDtos) == 0 {
		return nil, nil
	}

	icecreams, err := r.convert(icecreamsDtos)
	if err != nil {
		return nil, err
	}

	return icecreams, nil
}

func (r *IcecreamRepo) convert(dtos []dtos.Icecream) (icecreams []*domain.Icecream, err error) {
	for _, icecream := range dtos {
		icecreams = append(icecreams, &domain.Icecream{
			ProductID:             strconv.Itoa(icecream.ProductId),
			Name:                  icecream.Name,
			Description:           icecream.Description.String,
			Story:                 icecream.Story.String,
			ImageClosed:           icecream.ImageClosed.String,
			ImageOpen:             icecream.ImageOpen.String,
			AllergyInfo:           icecream.AllergyInfo.String,
			DietaryCertifications: icecream.DietaryCertifications.String,
		})
	}
	return icecreams, nil
}
