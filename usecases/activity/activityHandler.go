package activity

import (
	"encoding/json"
	"fit-byte/models"
	"fit-byte/types"
	"fit-byte/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type AcitivityHandler struct {
	activityService ActivityService
}

func NewActivityHandler(activityService ActivityService) AcitivityHandler {
	return AcitivityHandler{activityService}
}

func (h *AcitivityHandler) HandleCreateActivity(w http.ResponseWriter, r *http.Request) error {
	payload := struct {
		ActivityType      string    `json:"activityType" validate:"required,oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"`
		DoneAt            time.Time `json:"doneAt" validate:"required"`
		DurationInMinutes int       `json:"durationInMinutes" validate:"required,min=1"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return models.NewError(http.StatusBadRequest, err.Error())
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErr := fmt.Errorf("validation for '%s' failed", err.Field())
			return models.NewError(http.StatusBadRequest, validationErr.Error())
		}
	}

	newActivity, err := h.activityService.CreateActivity(models.Activity{
		ActivityType:      payload.ActivityType,
		DoneAt:            payload.DoneAt,
		DurationInMinutes: payload.DurationInMinutes,
	})
	if err != nil {
		return err
	}

	res := struct {
		ActivityId        string    `json:"activityId"`
		ActivityType      string    `json:"activityType"`
		DoneAt            time.Time `json:"doneAt"`
		DurationInMinutes int       `json:"durationInMinutes"`
		CaloriesBurned    int       `json:"caloriesBurned"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
	}{
		ActivityId:        newActivity.Id,
		ActivityType:      newActivity.ActivityType,
		DoneAt:            newActivity.DoneAt,
		DurationInMinutes: newActivity.DurationInMinutes,
		CaloriesBurned:    newActivity.CaloriesBurned,
		CreatedAt:         newActivity.CreatedAt,
		UpdatedAt:         newActivity.UpdatedAt,
	}
	utils.SetJsonResponse(w, http.StatusCreated, res)

	return nil
}

func (h *AcitivityHandler) HandleUpdateActivity(w http.ResponseWriter, r *http.Request) error {
	activityId := r.PathValue("activityId")
	payload := types.UpdateActivityPayload{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return models.NewError(http.StatusBadRequest, err.Error())
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErr := fmt.Errorf("validation for '%s' failed", err.Field())
			return models.NewError(http.StatusBadRequest, validationErr.Error())
		}
	}

	activity, err := h.activityService.UpdateActivity(activityId, payload)
	if err != nil {
		return err
	}

	res := struct {
		ActivityId        string    `json:"activityId"`
		ActivityType      string    `json:"activityType"`
		DoneAt            time.Time `json:"doneAt"`
		DurationInMinutes int       `json:"durationInMinutes"`
		CaloriesBurned    int       `json:"caloriesBurned"`
	}{
		ActivityId:        activity.Id,
		ActivityType:      activity.ActivityType,
		DoneAt:            activity.DoneAt,
		DurationInMinutes: activity.DurationInMinutes,
		CaloriesBurned:    activity.CaloriesBurned,
	}
	utils.SetJsonResponse(w, http.StatusOK, res)

	return nil
}