package user

import (
	"fit-byte/models"
	"fit-byte/utils"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return UserHandler{userService}
}

func (h *UserHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) error {
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
