package domain

// TODO Update & Delete

type Icecreamer interface {
	Create(icecream Icecream) (int, error)
	Creates(icecreams []Icecream) ([]int, error)
	Read(ids []int) ([]*Icecream, error)
}

type Ingredienter interface {
	Create(ingredient Ingredient) (int, error)
	Creates(ingredients Ingredients) ([]int, error)
	Read(icecreamProductId int) (Ingredients, error)
	Reads(icecreamProductIds []int) ([]Ingredients, error)
	ReadAll() (Ingredients, error)
}

type SourcingValuer interface {
	Create(sourcingValue SourcingValue) (int, error)
	Creates(sourcingValues SourcingValues) ([]int, error)
	Read(icecreamProductId int) (SourcingValues, error)
	Reads(icecreamProductIds []int) ([]SourcingValues, error)
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
