package models

import "time"

type Activity struct {
	Id                string    `json:"activityId" db:"id"`
	ActivityType      string    `json:"activityType" db:"activity_type"`
	DoneAt            time.Time `json:"doneAt" db:"done_at"`
	DurationInMinutes int       `json:"durationInMinutes" db:"duration_in_minutes"`
	CaloriesBurned    int       `json:"caloriesBurned" db:"calories_burned"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}
