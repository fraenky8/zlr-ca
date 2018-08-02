package dtos

type IcecreamHasSourcingValues struct {
	IcecreamProductId int64 `db:"IcecreamProductId"`
	SourcingValuesId  int64 `db:"SourcingValuesId"`
}
