package mock

type IcecreamHasIngredientsService struct {
	CreateFn      func(icecreamProductId int, ingredientIds []int) error
	CreateInvoked bool
}

func (s *IcecreamHasIngredientsService) Create(icecreamProductId int, ingredientIds []int) error {
	s.CreateInvoked = true
	return s.CreateFn(icecreamProductId, ingredientIds)
}
