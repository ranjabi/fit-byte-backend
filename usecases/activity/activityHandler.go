package activity

import (
	"encoding/json"
	"fit-byte/models"
	"fit-byte/types"
	"fit-byte/utils"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
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

func IsISO8601Date(fl validator.FieldLevel) bool {
	ISO8601DateRegexString := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
	ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
  return ISO8601DateRegex.MatchString(fl.Field().String())
}

func (h *AcitivityHandler) HandleGetAllActivities(w http.ResponseWriter, r *http.Request) error {
	
	validate := validator.New()
	validate.RegisterValidation("ISO8601date", IsISO8601Date)

	params := r.URL.Query()
	limitStr := params.Get("limit")
	offsetStr := params.Get("offset")
	activityTypeRaw := params.Get("activityType")
	doneAtFromRaw := params.Get("doneAtFrom")
	doneAtToRaw := params.Get("doneAtTo")
	caloriesBurnedMinRaw := params.Get("caloriesBurnedMin")
	caloriesBurnedMaxRaw := params.Get("caloriesBurnedMax")
	limit := 5
	offset := 0
	var activityType *string = &activityTypeRaw
	var doneAtFrom *string = &doneAtFromRaw
	var doneAtTo *string = &doneAtToRaw
	var caloriesBurnedMin *string = &caloriesBurnedMinRaw
	var caloriesBurnedMax *string = &caloriesBurnedMaxRaw
	
	if limitStr != "" {
		limitTemp, err := strconv.Atoi(limitStr)
		if err != nil {
			return models.NewError(http.StatusBadRequest, err.Error())
		}
		if limitTemp >= 0 {
			limit = limitTemp
		}
	}
	if offsetStr != "" {
		offsetTemp, err := strconv.Atoi(offsetStr)
		if err != nil {
			return models.NewError(http.StatusBadRequest, err.Error())
		}
		if offsetTemp >= 0 {
			offset = offsetTemp
		}
	}
	if err := validate.Var(activityType, "oneof=Walking Yoga Stretching Cycling Swimming Dancing Hiking Running HIIT JumpRope"); err != nil {
		// fmt.Println("----- SKIP activityType VALIDATION")
		activityType = nil
	}
	if err := validate.Var(doneAtFrom, "ISO8601date"); err != nil {
		// fmt.Println("----- SKIP doneAtFrom VALIDATION")
		doneAtFrom = nil
	}
	if err := validate.Var(doneAtTo, "ISO8601date"); err != nil {
		// fmt.Println("----- SKIP doneAtTo VALIDATION")
		doneAtTo = nil
	}
	if err := validate.Var(caloriesBurnedMin, "numeric"); err != nil {
		// fmt.Println("----- SKIP caloriesBurnedMin VALIDATION")
		caloriesBurnedMin = nil
	}
	if err := validate.Var(caloriesBurnedMax, "numeric"); err != nil {
		// fmt.Println("----- SKIP caloriesBurnedMax VALIDATION")
		caloriesBurnedMax = nil
	}

	activities, err := h.activityService.GetAllActivities(offset, limit, activityType, doneAtFrom, doneAtTo, caloriesBurnedMin, caloriesBurnedMax)
	if err != nil {
		return err
	}

	utils.SetJsonResponse(w, http.StatusOK, activities)

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

func (h *AcitivityHandler) HandleDeleteActivity(w http.ResponseWriter, r *http.Request) error {
	activityId := r.PathValue("activityId")

	err := h.activityService.DeleteActivity(activityId)
	if err != nil {
		return err
	}

	w.Write([]byte(""))

	return nil
}