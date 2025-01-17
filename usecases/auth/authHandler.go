package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"fit-byte/models"
	"fit-byte/utils"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) AuthHandler {
	return AuthHandler{authService}
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	payload := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=32"`
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

	newUser, err := h.authService.CreateUser(models.User{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		return err
	}

	res := struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		Email: newUser.Email,
		Token: newUser.Token,
	}
	utils.SetJsonResponse(w, http.StatusOK, res)

	return nil
}
