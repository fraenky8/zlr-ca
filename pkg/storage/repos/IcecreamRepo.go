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
	db   storage.Database
	repo Service
}

func NewIcecreamRepo(db storage.Database) *IcecreamRepo {

	repo := &IcecreamRepo{}

	service := Service{
		IcecreamService:                  repo,
		IngredientService:                NewIngredientsRepo(db),
		SourcingValueService:             NewSourcingValuesRepo(db),
		IcecreamHasIngredientsService:    NewIcecreamHasIngredientsRepo(db),
		IcecreamHasSourcingValuesService: NewIcecreamHasSourcingValuesRepo(db),
	}

	repo.db = db
	repo.repo = service

	return repo
}

func (r *IcecreamRepo) Creates(icecreams []*domain.Icecream) ([]int, error) {

	stmt, err := r.db.DB().Preparex(fmt.Sprintf(`
		INSERT INTO %s.icecream
  			(product_id, name, description, story, image_open, image_closed, allergy_info, dietary_certifications)
		VALUES
  			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING product_id
	`, r.db.Config().Schema))

	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}

	var ids []int
	for _, icecream := range icecreams {

		var productId int
		err = stmt.Get(&productId,
			icecream.ProductID, icecream.Name, icecream.Description, icecream.Story,
			icecream.ImageOpen, icecream.ImageClosed, icecream.AllergyInfo, icecream.DietaryCertifications,
		)
		if err != nil {
			return nil, fmt.Errorf("could not create icecream: %v", err)
		}

		ids, err := r.repo.IngredientService.Creates(icecream.Ingredients)
		if err != nil {
			return nil, err
		}

		err = r.repo.IcecreamHasIngredientsService.Create(productId, ids)
		if err != nil {
			return nil, err
		}

		ids, err = r.repo.SourcingValueService.Creates(icecream.SourcingValues)
		if err != nil {
			return nil, err
		}

		err = r.repo.IcecreamHasSourcingValuesService.Create(productId, ids)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		ids = append(ids, productId)
	}

	return ids, nil
}

func (r *IcecreamRepo) Reads(ids []int) ([]*domain.Icecream, error) {

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
	`, r.db.Config().Schema), ids)

	if err != nil {
		return nil, err
	}

	// http://jmoiron.github.io/sqlx/#inQueries
	// sqlx.In returns queries with the `?` bindvar, we can rebind it for our backend
	// here: ? to $#
	query = r.db.DB().Rebind(query)

	var icecreamsDtos []dtos.Icecream
	if err = r.db.DB().Select(&icecreamsDtos, query, args...); err != nil {
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

func (r *IcecreamRepo) Updates(icecreams []*domain.Icecream) (err error) {

	tx := r.db.DB().MustBegin()

	stmt, err := tx.Preparex(fmt.Sprintf(`
		UPDATE %s.icecream SET 
		  name = $1,
		  description = $2,
		  story = $3,
		  image_open = $4,
		  image_closed = $5,
		  allergy_info = $6,
		  dietary_certifications = $7
		WHERE product_id = $8
	`, r.db.Config().Schema))

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not prepare statement: %v", err)
	}

	for _, icecream := range icecreams {
		result, err := stmt.Exec(
			icecream.Name, icecream.Description,
			icecream.Story, icecream.ImageOpen,
			icecream.ImageClosed, icecream.AllergyInfo,
			icecream.DietaryCertifications, icecream.ProductID,
		)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("could not update icecream with productID = %s: %v", icecream.ProductID, err)
		}

		affectedRows, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("could not update icecream with productID = %s: %v", icecream.ProductID, err)
		}

		if affectedRows == 0 {
			tx.Rollback()
			return fmt.Errorf("icecream with productID = %s does not exist", icecream.ProductID)
		}
	}

	return tx.Commit()
}

func (r *IcecreamRepo) Deletes(ids []int) (err error) {

	tx := r.db.DB().MustBegin()

	stmt, err := tx.Preparex(fmt.Sprintf(`
		DELETE FROM %s.icecream
		WHERE product_id = $1
	`, r.db.Config().Schema))

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not prepare statement: %v", err)
	}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("could not update icecream with productID = %d: %v", id, err)
		}

		affectedRows, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("could not update icecream with productID = %d: %v", id, err)
		}

		if affectedRows == 0 {
			tx.Rollback()
			return fmt.Errorf("icecream with productID = %d does not exist", id)
		}
	}

	return tx.Commit()
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
