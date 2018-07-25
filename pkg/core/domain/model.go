package domain

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
