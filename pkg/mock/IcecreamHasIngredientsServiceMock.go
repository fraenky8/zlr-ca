package mock

type IcecreamHasIngredientsService struct {
	CreateFn      func(icecreamProductId int64, ingredientIds []int64) error
	CreateInvoked bool
}

func (s *IcecreamHasIngredientsService) Create(icecreamProductId int64, ingredientIds []int64) error {
	s.CreateInvoked = true
	return s.CreateFn(icecreamProductId, ingredientIds)
}
