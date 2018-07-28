package mock

import "github.com/fraenky8/zlr-ca/pkg/domain"

type IcecreamService struct {
	CreatesFn      func(icecreams []*domain.Icecream) ([]int, error)
	CreatesInvoked bool

	ReadsFn      func(ids []int) ([]*domain.Icecream, error)
	ReadsInvoked bool

	UpdatesFn      func(icecreams []*domain.Icecream) error
	UpdatesInvoked bool

	DeletesFn      func(ids []int) error
	DeletesInvoked bool
}

func (s *IcecreamService) Creates(icecreams []*domain.Icecream) ([]int, error) {
	s.CreatesInvoked = true
	return s.CreatesFn(icecreams)
}

func (s *IcecreamService) Reads(ids []int) ([]*domain.Icecream, error) {
	s.ReadsInvoked = true
	return s.ReadsFn(ids)
}

func (s *IcecreamService) Updates(icecreams []*domain.Icecream) error {
	s.UpdatesInvoked = true
	return s.UpdatesFn(icecreams)
}

func (s *IcecreamService) Deletes(ids []int) error {
	s.DeletesInvoked = true
	return s.DeletesFn(ids)
}
