package dtos

type IcecreamHasIngredients struct {
	IcecreamProductId int64 `db:"IcecreamProductId"`
	IngredientsId     int64 `db:"IngredientsId"`
}
