package mock

import "github.com/fraenky8/zlr-ca/pkg/domain"

type SourcingValueService struct {
	CreatesFn      func(sourcingValues domain.SourcingValues) ([]int, error)
	CreatesInvoked bool

	ReadFn      func(icecreamProductId int) (domain.SourcingValues, error)
	ReadInvoked bool

	ReadsFn      func(icecreamProductIds []int) ([]domain.SourcingValues, error)
	ReadsInvoked bool

	ReadAllFn      func() (domain.SourcingValues, error)
	ReadAllInvoked bool

	DeletesFn      func(icecreamProductIds []int) error
	DeletesInvoked bool
}

func (s *SourcingValueService) Creates(sourcingValues domain.SourcingValues) ([]int, error) {
	s.CreatesInvoked = true
	return s.CreatesFn(sourcingValues)
}

func (s *SourcingValueService) Reads(icecreamProductIds []int) ([]domain.SourcingValues, error) {
	s.ReadsInvoked = true
	return s.ReadsFn(icecreamProductIds)
}

func (s *SourcingValueService) Read(icecreamProductIds int) (domain.SourcingValues, error) {
	s.ReadsInvoked = true
	return s.ReadFn(icecreamProductIds)
}

func (s *SourcingValueService) ReadAll() (domain.SourcingValues, error) {
	s.ReadsInvoked = true
	return s.ReadAllFn()
}

func (s *SourcingValueService) Deletes(icecreamProductIds []int) error {
	s.DeletesInvoked = true
	return s.DeletesFn(icecreamProductIds)
}
