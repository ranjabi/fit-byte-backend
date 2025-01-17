package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"

	"fit-byte/constants"
	"fit-byte/models"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), constants.SALT_ROUND)
	return string(bytes), err
}

func CreateClaims(user *models.User) (string, error) {
	tokenAuth := jwtauth.New(constants.HASH_ALG, []byte(constants.JWT_SECRET), nil)
	claims := map[string]any{
		"user_id":    user.Id,
		"user_email": user.Email,
	}
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetJsonResponse(w http.ResponseWriter, statusCode int, response any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}

	return nil
}

func AppHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if err, ok := err.(*models.AppError); ok {
				if err.Code != 0 {
					http.Error(w, err.Error(), err.Code)
					return
				}
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}