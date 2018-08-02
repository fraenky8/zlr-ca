package dtos

import (
	"database/sql"
)

type Icecream struct {
	ProductId             int64          `db:"product_id"`
	Name                  string         `db:"name"`
	Description           sql.NullString `db:"description"`
	Story                 sql.NullString `db:"story"`
	ImageOpen             sql.NullString `db:"image_open"`
	ImageClosed           sql.NullString `db:"image_closed"`
	AllergyInfo           sql.NullString `db:"allergy_info"`
	DietaryCertifications sql.NullString `db:"dietary_certifications"`
}
