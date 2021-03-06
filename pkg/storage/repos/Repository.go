package repos

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/domain"
	"github.com/fraenky8/zlr-ca/pkg/storage"
)

type Repository struct {
	IcecreamService                  domain.IcecreamService
	IngredientService                domain.IngredientService
	SourcingValueService             domain.SourcingValueService
	IcecreamHasIngredientsService    domain.IcecreamHasIngredientsService
	IcecreamHasSourcingValuesService domain.IcecreamHasSourcingValuesService
}

func NewRepository(db storage.Database) (*Repository, error) {
	s := &Repository{
		IcecreamService:                  NewIcecreamRepo(db),
		IngredientService:                NewIngredientsRepo(db),
		SourcingValueService:             NewSourcingValuesRepo(db),
		IcecreamHasIngredientsService:    NewIcecreamHasIngredientsRepo(db),
		IcecreamHasSourcingValuesService: NewIcecreamHasSourcingValuesRepo(db),
	}

	if err := s.Verify(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Repository) Verify() error {
	if s.IcecreamService == nil {
		return fmt.Errorf("no IcecreamService given")
	}
	if s.IngredientService == nil {
		return fmt.Errorf("no IngredientService given")
	}
	if s.SourcingValueService == nil {
		return fmt.Errorf("no SourcingValueService given")
	}
	if s.IcecreamHasIngredientsService == nil {
		return fmt.Errorf("no IcecreamHasIngredientsService given")
	}
	if s.IcecreamHasSourcingValuesService == nil {
		return fmt.Errorf("no IcecreamHasSourcingValuesService given")
	}
	return nil
}
