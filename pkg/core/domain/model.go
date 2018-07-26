package domain

import "github.com/fraenky8/zlr-ca/pkg/infrastructure/storage/dtos"

type Icecreamer interface {
	Create(icecream Icecream) (int64, error)
	Creates(icecreams []Icecream) ([]int64, error)
	Read(id int64) (*Icecream, error)
	Reads(id []int64) ([]*Icecream, error)
}

type Converter interface {
	Convert(icecream *dtos.Icecream) (*Icecream, error)
}

type Ingredient string
type Ingredients []Ingredient

type SourcingValue string
type SourcingValues []SourcingValue

type Icecream struct {
	ProductID             string `json:"productId"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Story                 string `json:"story"`
	ImageClosed           string `json:"image_closed"`
	ImageOpen             string `json:"image_open"`
	AllergyInfo           string `json:"allergy_info"`
	DietaryCertifications string `json:"dietary_certifications"`
	SourcingValues        `json:"sourcing_values"`
	Ingredients           `json:"ingredients"`
}
