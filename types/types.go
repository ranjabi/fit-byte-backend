package types

import (
	"encoding/json"
	"time"
)

type UpdateUserPayload struct {
	Preference  *string         `json:"preference" db:"preference" validate:"oneof=CARDIO WEIGHT"`
	WeightUnit  *string         `json:"weightUnit" db:"weight_unit" validate:"oneof=KG LBS"`
	HeightUnit  *string         `json:"heightUnit" db:"height_unit" validate:"oneof=CM INCH"`
	Weight      *int            `json:"weight" db:"weight" validate:"gte=10,lte=1000"`
	Height      *int            `json:"height" db:"height" validate:"gte=3,lte=250"`
	NameRaw     json.RawMessage `json:"name,omitempty"`
	Name        *string         `db:"name" validate:"omitempty,min=2,max=60"`
	ImageUriRaw json.RawMessage `json:"imageUri,omitempty"`
	ImageUri    *string         `db:"image_uri" validate:"omitempty,uri"`
}

type UpdateActivityPayload struct {
	ActivityTypeRaw   json.RawMessage `json:"activityType,omitempty"`
	ActivityType      *string         `db:"activity_type" validate:"omitempty,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"`
	DoneAtRaw         json.RawMessage `json:"doneAt,omitempty"`
	DoneAt            *time.Time      `db:"done_at" validate:"omitempty"`
	DurationInMinutes *int            `json:"durationInMinutes,omitempty" db:"duration_in_minutes" validate:"omitempty,min=1"`
	CaloriesBurned    *int             `json:"-" db:"calories_burned"`
}
