package dtos

import (
	"database/sql"
)

type Icecream struct {
	ProductId             int            `db:"ProductId"`
	Name                  string         `db:"Name"`
	Description           string         `db:"Description"`
	Story                 sql.NullString `db:"Story"`
	ImageOpen             sql.NullString `db:"ImageOpen"`
	ImageClosed           sql.NullString `db:"ImageClosed"`
	AllergyInfo           sql.NullString `db:"AllergyInfo"`
	DietaryCertifications sql.NullString `db:"DietaryCertifications"`
}
