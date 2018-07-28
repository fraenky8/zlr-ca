package domain

import (
	"fmt"
	"strings"
)

// TODO Update & Delete

type IcecreamService interface {
	Create(icecream Icecream) (int, error)
	Creates(icecreams []Icecream) ([]int, error)
	Read(ids []int) ([]*Icecream, error)
}

type IngredientService interface {
	Create(ingredient Ingredient) (int, error)
	Creates(ingredients Ingredients) ([]int, error)
	Read(icecreamProductId int) (Ingredients, error)
	Reads(icecreamProductIds []int) ([]Ingredients, error)
	ReadAll() (Ingredients, error)
}

type SourcingValueService interface {
	Create(sourcingValue SourcingValue) (int, error)
	Creates(sourcingValues SourcingValues) ([]int, error)
	Read(icecreamProductId int) (SourcingValues, error)
	Reads(icecreamProductIds []int) ([]SourcingValues, error)
	ReadAll() (SourcingValues, error)
}

type Ingredient string
type Ingredients []Ingredient

func (i Ingredient) Verify() error {
	if strings.TrimSpace(string(i)) == "" {
		return fmt.Errorf("missing valid ingredient name")
	}
	return nil
}

func (is Ingredients) Verify() error {
	for _, i := range is {
		if err := i.Verify(); err != nil {
			return err
		}
	}
	return nil
}

type SourcingValue string
type SourcingValues []SourcingValue

func (s SourcingValue) Verify() error {
	if strings.TrimSpace(string(s)) == "" {
		return fmt.Errorf("missing valid sourcing value description")
	}
	return nil
}

func (sv SourcingValues) Verify() error {
	for _, s := range sv {
		if err := s.Verify(); err != nil {
			return err
		}
	}
	return nil
}

type Icecream struct {
	ProductID             string `json:"productId"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Story                 string `json:"story"`
	ImageClosed           string `json:"image_closed"`
	ImageOpen             string `json:"image_open"`
	AllergyInfo           string `json:"allergy_info"`
	DietaryCertifications string `json:"dietary_certifications"`
	SourcingValues        `json:"sourcing_values,omitempty"`
	Ingredients           `json:"ingredients,omitempty"`
}

func (i Icecream) Verify() error {
	if strings.TrimSpace(i.ProductID) == "" {
		return fmt.Errorf("missing valid product id")
	}
	if strings.TrimSpace(i.Name) == "" {
		return fmt.Errorf("missing valid name")
	}
	return nil
}
