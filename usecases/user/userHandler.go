package user

import (
	"encoding/json"
	"fit-byte/models"
	"fit-byte/types"
	"fit-byte/utils"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return UserHandler{userService}
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return models.NewError(http.StatusInternalServerError, err.Error())
	}

	user, err := h.userService.FindById(claims["userId"].(string))
	if err != nil {
		return models.NewError(http.StatusInternalServerError, err.Error())
	}

	res := struct {
		Preference pgtype.Text `json:"preference"`
		WeightUnit pgtype.Text `json:"weightUnit"`
		HeightUnit pgtype.Text `json:"heightUnit"`
		Weight     pgtype.Int4 `json:"weight"`
		Height     pgtype.Int4 `json:"height"`
		Email      string      `json:"email"`
		Name       pgtype.Text `json:"name"`
		ImageUri   pgtype.Text `json:"imageuri"`
	}{
		Preference: user.Preference,
		WeightUnit: user.WeightUnit,
		HeightUnit: user.HeightUnit,
		Weight:     user.Weight,
		Height:     user.Height,
		Email:      user.Email,
		Name:       user.Name,
		ImageUri:   user.ImageUri,
	}
	utils.SetJsonResponse(w, http.StatusOK, res)

	return nil
}

func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	payload := types.UpdateUserPayload{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return models.NewError(http.StatusBadRequest, err.Error())
	}
	if !utils.IsValidURI(*payload.ImageUri) {
		return models.NewError(http.StatusBadRequest, "Invalid uri")
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErr := fmt.Errorf("validation for '%s' failed", err.Field())
			return models.NewError(http.StatusBadRequest, validationErr.Error())
		}
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return models.NewError(http.StatusInternalServerError, err.Error())
	}
	userId := claims["userId"].(string)
	user, err := h.userService.PartialUpdate(userId, payload)
	if err != nil {
		return err
	}

	res := struct {
		Preference pgtype.Text `json:"preference"`
		WeightUnit pgtype.Text `json:"weightUnit"`
		HeightUnit pgtype.Text `json:"heightUnit"`
		Weight     pgtype.Int4 `json:"weight"`
		Height     pgtype.Int4 `json:"height"`
		Name       pgtype.Text `json:"name"`
		ImageUri   pgtype.Text `json:"imageuri"`
	}{
		Preference: user.Preference,
		WeightUnit: user.WeightUnit,
		HeightUnit: user.HeightUnit,
		Weight:     user.Weight,
		Height:     user.Height,
		Name:       user.Name,
		ImageUri:   user.ImageUri,
	}
	utils.SetJsonResponse(w, http.StatusOK, res)

	return nil
}
