package domain

type Icecreamer interface {
	Create(icecream Icecream) (int64, error)
	Creates(icecreams []Icecream) ([]int64, error)
	Read(id int64) (*Icecream, error)
	Reads(id []int64) ([]*Icecream, error)
	// TODO Update & Delete
}

type Ingredienter interface {
	Create(ingredient Ingredient)
	Creates(ingredients Ingredients) ([]int64, error)
	Read(icecreamProductId int64) (*Ingredients, error)
}

type SourcingValuer interface {
	Create(sourcingValue SourcingValue)
	Creates(sourcingValues SourcingValues) ([]int64, error)
	Read(icecreamProductId int64) (*SourcingValues, error)
}

// type IcecreamConverter interface {
// 	Convert(icecream *dtos.Icecream) (*Icecream, error)
// }
// type IngredientConverter interface {
// 	Convert(ingredient []*dtos.Ingredients) (*Ingredients, error)
// }
// type SourcingValueConverter interface {
// 	Convert(sourcingValues []*dtos.SourcingValues) (domain.SourcingValues, error)
// }

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
