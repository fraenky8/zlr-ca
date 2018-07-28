package dtos

type IcecreamHasIngredients struct {
	IcecreamProductId int `db:"IcecreamProductId"`
	IngredientsId     int `db:"IngredientsId"`
}
