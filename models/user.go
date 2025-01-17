package models

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	Id         string      `json:"id" db:"id"`
	Email      string      `json:"email" db:"email"`
	Password   string      `json:"-" db:"password"`
	Token      string      `json:"token" db:"-"`
	Preference pgtype.Text `json:"preference" db:"preference"`
	WeightUnit pgtype.Text `json:"WeightUnit" db:"weight_unit"`
	HeightUnit pgtype.Text `json:"HeightUnit" db:"height_unit"`
	Weight     pgtype.Int4 `json:"weight" db:"weight"`
	Height     pgtype.Int4 `json:"height" db:"height"`
	Name       pgtype.Text `json:"name" db:"name"`
	ImageUri   pgtype.Text `json:"imageUri" db:"image_uri"`
}
