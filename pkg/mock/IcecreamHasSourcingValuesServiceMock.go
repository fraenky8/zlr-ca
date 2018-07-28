package mock

type IcecreamHasSourcingValuesService struct {
	CreateFn      func(icecreamProductId int, ingredientIds []int) error
	CreateInvoked bool
}

func (s *IcecreamHasSourcingValuesService) Create(icecreamProductId int, sourcingValueIds []int) error {
	s.CreateInvoked = true
	return s.CreateFn(icecreamProductId, sourcingValueIds)
}
