package types

type UpdateUserPayload struct {
	Preference *string `json:"preference" db:"preference" validate:"oneof=CARDIO WEIGHT"`
	WeightUnit *string `json:"weightUnit" db:"weight_unit" validate:"oneof=KG LBS"`
	HeightUnit *string `json:"heightUnit" db:"height_unit" validate:"oneof=CM INCH"`
	Weight     *int    `json:"weight" db:"weight" validate:"gte=10,lte=1000"`
	Height     *int    `json:"height" db:"height" validate:"gte=3,lte=250"`
	Name       *string `json:"name,omitempty" db:"name" validate:"omitempty,min=2,max=60"`
	ImageUri   *string `json:"imageuri,omitempty" db:"image_uri" validate:"omitempty,uri"`
}
