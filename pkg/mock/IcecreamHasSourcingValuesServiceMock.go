package mock

type IcecreamHasSourcingValuesService struct {
	CreateFn      func(icecreamProductId int64, ingredientIds []int64) error
	CreateInvoked bool
}

func (s *IcecreamHasSourcingValuesService) Create(icecreamProductId int64, sourcingValueIds []int64) error {
	s.CreateInvoked = true
	return s.CreateFn(icecreamProductId, sourcingValueIds)
}
