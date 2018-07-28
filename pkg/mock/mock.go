package mock

import "github.com/fraenky8/zlr-ca/pkg/domain"

type IcecreamService struct {
	CreateFn      func(icecream domain.Icecream) (int, error)
	CreateInvoked bool

	CreatesFn      func(icecreams []domain.Icecream) ([]int, error)
	CreatesInvoked bool

	ReadFn      func(ids []int) ([]*domain.Icecream, error)
	ReadInvoked bool
}

func (s *IcecreamService) Create(icecream domain.Icecream) (int, error) {
	s.CreateInvoked = true
	return s.CreateFn(icecream)
}

func (s *IcecreamService) Creates(icecreams []domain.Icecream) ([]int, error) {
	s.CreatesInvoked = true
	return s.CreatesFn(icecreams)
}

func (s *IcecreamService) Read(ids []int) ([]*domain.Icecream, error) {
	s.ReadInvoked = true
	return s.ReadFn(ids)
}
