package mock

import "github.com/fraenky8/zlr-ca/pkg/domain"

type IngredientService struct {
	CreatesFn      func(ingredients domain.Ingredients) ([]int, error)
	CreatesInvoked bool

	ReadFn      func(icecreamProductId int) (domain.Ingredients, error)
	ReadInvoked bool

	ReadsFn      func(icecreamProductIds []int) ([]domain.Ingredients, error)
	ReadsInvoked bool

	ReadAllFn      func() (domain.Ingredients, error)
	ReadAllInvoked bool
}

func (s *IngredientService) Creates(ingredients domain.Ingredients) ([]int, error) {
	s.CreatesInvoked = true
	return s.CreatesFn(ingredients)
}

func (s *IngredientService) Reads(icecreamProductIds []int) ([]domain.Ingredients, error) {
	s.ReadsInvoked = true
	return s.ReadsFn(icecreamProductIds)
}

func (s *IngredientService) Read(icecreamProductIds int) (domain.Ingredients, error) {
	s.ReadsInvoked = true
	return s.ReadFn(icecreamProductIds)
}

func (s *IngredientService) ReadAll() (domain.Ingredients, error) {
	s.ReadsInvoked = true
	return s.ReadAllFn()
}
