package domain

// TODO Update & Delete

type Icecreamer interface {
	Create(icecream Icecream) (int64, error)
	Creates(icecreams []Icecream) ([]int64, error)
	Read(ids []int64) ([]*Icecream, error)
}

type Ingredienter interface {
	Create(ingredient Ingredient) (int64, error)
	Creates(ingredients Ingredients) ([]int64, error)
	Read(icecreamProductId int64) (Ingredients, error)
	Reads(icecreamProductIds []int64) ([]Ingredients, error)
	ReadAll() (Ingredients, error)
}

type SourcingValuer interface {
	Create(sourcingValue SourcingValue) (int64, error)
	Creates(sourcingValues SourcingValues) ([]int64, error)
	Read(icecreamProductId int64) (SourcingValues, error)
	Reads(icecreamProductIds []int64) ([]SourcingValues, error)
	ReadAll() (SourcingValues, error)
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
	SourcingValues        `json:"sourcing_values,omitempty"`
	Ingredients           `json:"ingredients,omitempty"`
}
